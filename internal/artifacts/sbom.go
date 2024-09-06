package artifacts

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta/mta"
	"github.com/pkg/errors"
)

const (
	xml_type         = "xml"
	json_type        = "json"
	unsupport_type   = "unsupport_sbom_type"
	xml_suffix       = ".xml"
	json_suffix      = ".json"
	sbom_xml_suffix  = ".bom.xml"
	sbom_json_suffix = ".bom.json"
	cyclonedx_cli    = "cyclonedx"
)

type Bom struct {
	XMLName  xml.Name `xml:"bom"`
	Metadata Metadata `xml:"metadata"`
}

type Metadata struct {
	XMLName   xml.Name  `xml:"metadata"`
	Component Component `xml:"component"`
}

type Component struct {
	XMLName xml.Name `xml:"component"`
	BomRef  string   `xml:"bom-ref,attr"`
}

type Dependency struct {
	XMLName    xml.Name `xml:"dependency"`
	Ref        string   `xml:"ref,attr"`
	SubDepends []SubDep `xml:"dependency"`
}

type SubDep struct {
	XMLName xml.Name `xml:"dependency"`
	Ref     string   `xml:"ref,attr"`
}

// ExecuteProjectSBomGenerate - Execute MTA project SBOM generation
func ExecuteProjectSBomGenerate(source string, sbomFilePath string, wdGetter func() (string, error)) error {
	// (1) get loc object and mta object
	loc, err := dir.Location(source, "", dir.Dev, []string{}, wdGetter)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	mtaObj, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	// (2) if sbom file path is empty, default value is <MTA project path>/<MTA project id>.bom.xml
	if strings.TrimSpace(sbomFilePath) == "" {
		sbomFilePath = mtaObj.ID + sbom_xml_suffix
	}

	// (3) generate sbom
	err = executeSBomGenerate(loc, mtaObj, source, sbomFilePath)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	return nil
}

// ExecuteProjectBuildeSBomGenerate - Execute MTA project SBOM generation with Build process
func ExecuteProjectBuildeSBomGenerate(source string, sbomFilePath string, wdGetter func() (string, error)) error {
	// (1) if sbomFilePath is empty, do not need to generate sbom, return directly
	if strings.TrimSpace(sbomFilePath) == "" {
		return nil
	}

	// (2) get loc object and mta object
	loc, err := dir.Location(source, "", dir.Dev, []string{}, wdGetter)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	mtaObj, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	// (3) generate sbom
	err = executeSBomGenerate(loc, mtaObj, source, sbomFilePath)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	return nil
}

// prepareEnv - create sbom tmp dir and sbom target dir
// Notice: must not remove the sbom target dir, if the sbom-file is generated under project root, remove sbom path will delete app
func prepareEnv(sbomTmpDir, sbomPath string) error {
	err := dir.RemoveIfExist(sbomTmpDir)
	if err != nil {
		return err
	}
	err = dir.CreateDirIfNotExist(sbomTmpDir)
	if err != nil {
		return err
	}
	if len(sbomPath) > 0 {
		err = dir.CreateDirIfNotExist(sbomPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// cleanEnv - clean sbom tmp dir
func cleanEnv(sbomTmpDir string) error {
	err := dir.RemoveIfExist(sbomTmpDir)
	if err != nil {
		return err
	}
	return nil
}

// generateSBomFile - generate all modules sbom and merge in to one, then mv it to sbom target path
func generateSBomFile(loc *dir.Loc, mtaObj *mta.MTA,
	sbomPath, sbomName, sbomType, sbomSuffix, sbomTmpDir string) error {
	// (1) generation sbom for modules under sbom tmp dir
	err := generateSBomFiles(loc, mtaObj, sbomTmpDir, sbomType, sbomSuffix)
	if err != nil {
		return err
	}

	// (2) list all generated module sbom file in tmp
	sbomFileNames, err := listSBomFilesInTmpDir(sbomTmpDir, sbomSuffix)
	if err != nil {
		return err
	}
	// if sbom tmp dir is empty, maybe all modules are unknow builder or custom builders
	if len(sbomFileNames) == 0 {
		logs.Logger.Infof(genSBomEmptyMsg, sbomName)
		return nil
	}

	// (3) merge sbom files under sbom tmp dir
	sbomTmpName, err := mergeSBomFiles(loc, mtaObj, sbomTmpDir, sbomFileNames, sbomName, sbomType, sbomSuffix)
	if err != nil {
		return err
	}

	// (4) get module bom-ref info
	moduleBomRefs, err := getModuleBomRefs(sbomTmpDir, sbomFileNames)

	for _, bomRef := range moduleBomRefs {
		logs.Logger.Infof("moduleBomRef:%s", bomRef)
	}
	if err != nil {
		return err
	}

	// (4) instert xml attribute or xml node to bom->metadata
	err = updateSBomMetadataNode(mtaObj, sbomTmpDir, sbomTmpName, moduleBomRefs)
	if err != nil {
		return err
	}

	// (5) generate sbom target dir, mv merged sbom file to target dir
	err = moveSBomToTarget(sbomPath, sbomName, sbomTmpDir, sbomTmpName)
	if err != nil {
		return err
	}

	return nil
}

func getModuleBomRefs(sbomTmpDir string, sbomFileNames []string) ([]string, error) {
	bomRefMap := make(map[string]struct{})

	for _, fileName := range sbomFileNames {
		sbomfilepath := filepath.Join(sbomTmpDir, fileName)
		xmlFile, err := os.Open(sbomfilepath)
		if err != nil {
			return nil, err
		}
		defer xmlFile.Close()

		byteValue, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			return nil, err
		}

		var bom Bom
		if err := xml.Unmarshal(byteValue, &bom); err != nil {
			return nil, err
		}

		bomRefMap[bom.Metadata.Component.BomRef] = struct{}{}
	}

	var moduleBomRefs []string
	for bomRef := range bomRefMap {
		moduleBomRefs = append(moduleBomRefs, bomRef)
	}

	return moduleBomRefs, nil
}

func removeXmlns(attrs []xml.Attr) []xml.Attr {
	var result []xml.Attr
	for _, attr := range attrs {
		if attr.Name.Local != "xmlns" {
			result = append(result, attr)
		}
	}
	return result
}

func addBomrefAttribute(attributes []xml.Attr, purl string) []xml.Attr {
	purlAttr := xml.Attr{
		Name:  xml.Name{Local: "bom-ref"},
		Value: purl,
	}

	// Add bom-ref attribute to attributes list
	attributes = append(attributes, purlAttr)

	return attributes
}

func addXmlnsSchemaAttribute(attributes []xml.Attr, xmlnsSchema string) []xml.Attr {
	purlAttr := xml.Attr{
		Name:  xml.Name{Local: "xmlns"},
		Value: xmlnsSchema,
	}

	// Add bom-ref attribute to attributes list
	attributes = append(attributes, purlAttr)

	return attributes
}

func updateSBomMetadataNode(mtaObj *mta.MTA, sbomTmpDir, sbomTmpName string, moduleBomRefs []string) error {
	sbomfilepath := filepath.Join(sbomTmpDir, sbomTmpName)
	file, err := os.Open(sbomfilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var purl = "pkg:mta/" + mtaObj.ID + "@" + mtaObj.Version
	var xmlnsSchema = "http://cyclonedx.org/schema/bom/1.4"

	decoder := xml.NewDecoder(file)
	decoder.Strict = false

	var out bytes.Buffer
	encoder := xml.NewEncoder(&out)

	isInBom := false
	isInBomMetadata := false

	for {
		tok, err := decoder.RawToken()
		// tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch typedTok := tok.(type) {
		case xml.ProcInst:
			out.Write([]byte("<?" + string(typedTok.Target) + " " + string(typedTok.Inst) + "?>"))
		case xml.StartElement:
			// if xml node contains 'xmlns' attribute, remove the attribute
			typedTok.Attr = removeXmlns(typedTok.Attr)
			// if current node is <bom>
			if typedTok.Name.Local == "bom" {
				// 1. set isInBom = true
				isInBom = true
				// 2. add xmlns schema attribute to bom xml node
				typedTok.Attr = addXmlnsSchemaAttribute(typedTok.Attr, xmlnsSchema)
			}
			// if current node is bom->metadata
			if typedTok.Name.Local == "metadata" && isInBom {
				isInBomMetadata = true
				// 1. write bom->metadata xml node
				err := encoder.EncodeToken(typedTok)
				if err != nil {
					return err
				}
				// 2. add bom->meatadata->timestamp xml node
				encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "timestamp"}})
				encoder.EncodeToken(xml.CharData(time.Now().UTC().Format("2006-01-02T15:04:05Z")))
				encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "timestamp"}})
				break
			}

			// if current node is bom-metadata->component
			if typedTok.Name.Local == "component" && isInBom && isInBomMetadata {
				// 1. add bom-ref attribute to bom->metadata->component xml node
				typedTok.Attr = addBomrefAttribute(typedTok.Attr, purl)
				// 2. write bom->metadata->component xml node
				err := encoder.EncodeToken(typedTok)
				if err != nil {
					return err
				}
				// 3. add purl xml node to bom->metadata->component xml node
				encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "purl"}})
				encoder.EncodeToken(xml.CharData(purl))
				encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "purl"}})
				break
			}

			if typedTok.Name.Local == "dependencies" && isInBom {
				// 1. write bom->dependencies xml node
				err := encoder.EncodeToken(typedTok)
				if err != nil {
					return err
				}
				// 2. Todo: insert new dependency xml element
				dependency := Dependency{}
				dependency.Ref = purl
				for _, bomRefs := range moduleBomRefs {
					subDependency := SubDep{}
					subDependency.Ref = bomRefs
					dependency.SubDepends = append(dependency.SubDepends, subDependency)
				}
				encoder.Encode(dependency)
				break
			}

			// common xml node
			err := encoder.EncodeToken(typedTok)
			if err != nil {
				return err
			}
		case xml.CharData:
			err := encoder.EncodeToken(typedTok)
			if err != nil {
				return err
			}
		case xml.EndElement:
			if typedTok.Name.Local == "bom" {
				isInBom = false
			}
			if typedTok.Name.Local == "metadata" && isInBom {
				isInBomMetadata = false
			}

			err := encoder.EncodeToken(typedTok)
			if err != nil {
				return err
			}
		default:
			err := encoder.EncodeToken(typedTok)
			if err != nil {
				return err
			}
		}
	}
	encoder.Flush()
	content := out.Bytes()
	content = bytes.Replace(content, []byte("\ufeff"), []byte(""), -1)
	err = ioutil.WriteFile(sbomfilepath, content, 0644)

	if err != nil {
		return err
	}
	return nil
}

func executeSBomGenerate(loc *dir.Loc, mtaObj *mta.MTA, source string, sbomFilePath string) error {
	// start generate sbom file log
	logs.Logger.Info(genSBomFileStartMsg)

	// (1) parse sbomFilePath, if relative, it is relative path to project source
	// json type sbom file is not supported at present, if sbom file type is json, return not support error
	sbomPath, sbomName, sbomType, sbomSuffix := parseSBomFilePath(loc.GetSource(), sbomFilePath)
	// logs.Logger.Infof("source: %s; sbomFilePath: %s", loc.GetSource(), sbomFilePath)
	// logs.Logger.Infof("sbomPath: %s; sbomName: %s; sbomType: %s; sbomSuffix: %s", sbomPath, sbomName, sbomType, sbomSuffix)

	if sbomType == unsupport_type {
		return errors.Errorf(genSBomNotSupportedFileTypeMsg, sbomSuffix)
	}

	// (2) create sbom tmp dir and sbom target path
	sbomTmpDir := loc.GetSBomFileTmpDir(mtaObj)
	prepareErr := prepareEnv(sbomTmpDir, sbomPath)
	if prepareErr != nil {
		return prepareErr
	}

	// (3) generate sbom file
	genError := generateSBomFile(loc, mtaObj, sbomPath, sbomName, sbomType, sbomSuffix, sbomTmpDir)
	if genError != nil {
		cleanErr := cleanEnv(sbomTmpDir)
		if cleanErr != nil {
			logs.Logger.Error(cleanErr)
		}
		return genError
	}

	// (4) clean sbom tmp dir
	cleanErr := cleanEnv(sbomTmpDir)
	if cleanErr != nil {
		return cleanErr
	}

	// finish generate sbom file log
	logs.Logger.Infof(genSBomFileFinishedMsg, sbomName)

	return nil
}

// moveSBomToTarget - move sbom file from sbom tmp dir to target dir
func moveSBomToTarget(sbomPath string, sbomName string, sbomTmpDir string, sbomTmpName string) error {
	if len(sbomPath) > 0 {
		err := dir.CreateDirIfNotExist(sbomPath)
		if err != nil {
			return errors.Wrapf(err, createSBomTargetDirFailedMsg, sbomName)
		}
	}

	sourcesbomfilepath := filepath.Join(sbomTmpDir, sbomTmpName)
	targetsbomfilepath := filepath.Join(sbomPath, sbomName)

	err := os.Rename(sourcesbomfilepath, targetsbomfilepath)
	if err != nil {
		return errors.Wrapf(err, mvSBomToTargetDirFailedMsg, sourcesbomfilepath, targetsbomfilepath)
	}
	return nil
}

// listSBomFilesInTmpDir - list generated sbom files for modules
// if sbom tmp dir is empty, return empty array
func listSBomFilesInTmpDir(sbomTmpDir, sbomSuffix string) ([]string, error) {
	var sbomFileNames []string
	fileInfos, err := ioutil.ReadDir(sbomTmpDir)
	if err != nil {
		return sbomFileNames, err
	}

	for _, file := range fileInfos {
		fileName := file.Name()
		if !file.IsDir() && len(fileName) > 0 && strings.HasSuffix(fileName, sbomSuffix) {
			sbomFileNames = append(sbomFileNames, fileName)
		}
	}
	return sbomFileNames, nil
}

// mergeSBomFiles - merge sbom files of modules under sbom tmp dir
func mergeSBomFiles(loc *dir.Loc, mtaObj *mta.MTA, sbomTmpDir string, sbomFileNames []string, sbomName, sbomType, sbomSuffix string) (string, error) {
	curtime := time.Now().Format("20230328150313")

	var sbomTmpName string
	if strings.HasSuffix(sbomName, sbom_xml_suffix) {
		sbomTmpName = strings.TrimSuffix(sbomName, xml_suffix) + "_" + curtime + sbom_xml_suffix
	} else if strings.HasSuffix(sbomName, sbom_json_suffix) {
		sbomTmpName = strings.TrimSuffix(sbomName, json_suffix) + "_" + curtime + sbom_json_suffix
	} else {
		sbomTmpName = sbomName + "_" + curtime + sbom_xml_suffix
	}

	// get sbom file generate command
	sbomMergeCmds, err := commands.GetSBomsMergeCommand(loc, cyclonedx_cli, mtaObj, sbomTmpDir, sbomFileNames, sbomTmpName, sbomType, sbomSuffix)
	if err != nil {
		return "", err
	}

	// merging sbom file log
	logs.Logger.Infof(genSBomFileMergingMsg, sbomName)

	// exec sbom merge command
	err = executeSBomCommand(sbomMergeCmds)
	if err != nil {
		return "", err
	}
	return sbomTmpName, nil
}

// parseSBomFilePath - parse sbom file path parameter
// if sbom file path is a relative path, join source path
// only xml file format is supported at present;
func parseSBomFilePath(source string, sbomFilePath string) (string, string, string, string) {
	var sbomPath, sbomName, sbomType, sbomSuffix string
	if filepath.IsAbs(sbomFilePath) {
		sbomPath, sbomName = filepath.Split(sbomFilePath)
	} else {
		sbomPath, sbomName = filepath.Split(filepath.Join(source, sbomFilePath))
	}

	// if file suffix is .xml, or no file suffix, xml format type will be return
	// if file suffix is not .xml, unsupported file type will be return
	fileSuffix := filepath.Ext(sbomName)
	if fileSuffix == "" || strings.HasSuffix(sbomName, xml_suffix) {
		sbomType = xml_type
		sbomSuffix = sbom_xml_suffix
	} else {
		sbomType = unsupport_type
		sbomSuffix = fileSuffix
	}

	return sbomPath, sbomName, sbomType, sbomSuffix
}

func executeSBomCommand(sbomCmds [][]string) error {
	err := exec.ExecuteWithTimeout(sbomCmds, "", true)
	if err != nil {
		return err
	}
	return nil
}

// generateSBomFiles - loop all mta modules and generate sbom for each of then
// if module's builder is custom, skip it
func generateSBomFiles(loc *dir.Loc, mtaObj *mta.MTA, sBomFileTmpDir string, sbomType string, sbomSuffix string) error {
	// (1) sort module by dependency orders
	sortedModuleNames, err := buildops.GetModulesNames(mtaObj)
	if err != nil {
		return err
	}

	// (2) loop modules to generate sbom files
	curtime := time.Now().Format("20230328150313")
	for _, moduleName := range sortedModuleNames {
		// start generate module sbom log
		logs.Logger.Infof(genSBomForModuleStartMsg, moduleName)

		module, err := mtaObj.GetModuleByName(moduleName)
		if err != nil {
			return err
		}

		sbomFileName := moduleName + "_" + curtime
		sbomFileFullName := sbomFileName + sbomSuffix

		// get sbom file generate command
		sbomGenCmds, err := commands.GetModuleSBomGenCommands(loc, module, sbomFileName, sbomType, sbomSuffix)
		if err != nil {
			return err
		}
		// if sbomGenCmds is empty, module builder maybe "custom" or unknow builder, skip the module and continue
		if len(sbomGenCmds) == 0 {
			logs.Logger.Infof(genSBomSkipModuleMsg, moduleName)
			continue
		}

		// exec sbom generate command
		err = executeSBomCommand(sbomGenCmds)
		if err != nil {
			return err
		}

		// mv module sbom file to sbom temp dir
		modulePath := loc.GetSourceModuleDir(module.Path)
		sbomFileFoundPath, err := dir.FindFile(modulePath, sbomFileFullName)
		if err != nil {
			return err
		}
		sbomFileTargetPath := filepath.Join(sBomFileTmpDir, sbomFileFullName)
		err = os.Rename(sbomFileFoundPath, sbomFileTargetPath)
		if err != nil {
			return err
		}

		// finish generate module sbom log
		logs.Logger.Infof(genSBomForModuleFinishMsg, moduleName)
	}

	return nil
}

// ExecuteModuleSBomGenerate - Execute specified modules of MTA project SBOM generation
func ExecuteModuleSBomGenerate(source string, modulesNames []string, allDependencies bool, sBomFilePath string, wdGetter func() (string, error)) error {
	logs.Logger.Info("source: " + source)

	for _, moduleName := range modulesNames {
		logs.Logger.Info("module: " + moduleName)
	}

	logs.Logger.Info("allDependencies: " + strconv.FormatBool(allDependencies))
	logs.Logger.Info("sBomFilePath: " + sBomFilePath)

	message, err := version.GetVersionMessage()
	if err == nil {
		logs.Logger.Info(message)
	}
	return err
}
