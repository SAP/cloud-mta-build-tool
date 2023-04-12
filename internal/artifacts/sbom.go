package artifacts

import (
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
	xml_suffix       = ".xml"
	json_suffix      = ".json"
	sbom_xml_suffix  = ".bom.xml"
	sbom_json_suffix = ".bom.json"
	cyclonedx_cli    = "cyclonedx-cli"
)

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
	err = dir.CreateDirIfNotExist(sbomPath)
	if err != nil {
		return err
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

func generateSBomFile(loc *dir.Loc, mtaObj *mta.MTA, sbomPath, sbomName, sbomType, sbomSuffix, sbomTmpDir string) error {
	// (1) generation sbom for modules under sbom tmp dir
	err := generateSBomFiles(loc, mtaObj, sbomTmpDir, sbomType, sbomSuffix)
	if err != nil {
		return err
	}

	// (2) merge sbom files under sbom tmp dir
	sbomTmpName, err := mergeSBomFiles(loc, sbomTmpDir, sbomName, sbomSuffix)
	if err != nil {
		return err
	}

	// (3) generate sbom target dir, mv merged sbom file to target dir
	err = moveSBomToTarget(sbomPath, sbomName, sbomTmpDir, sbomTmpName)
	if err != nil {
		return err
	}

	return nil
}

func executeSBomGenerate(loc *dir.Loc, mtaObj *mta.MTA, source string, sbomFilePath string) error {
	// (1) parse sbomFilePath, if relative, it is relative path to project source
	sbomPath, sbomName, sbomType, sbomSuffix := parseSBomFilePath(loc.GetSource(), sbomFilePath)

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

	return nil
}

// moveSBomToTarget - move sbom file from sbom tmp dir to target dir
func moveSBomToTarget(sbomPath string, sbomName string, sbomTmpDir string, sbomTmpName string) error {
	err := dir.CreateDirIfNotExist(sbomPath)
	if err != nil {
		return errors.Wrapf(err, createSBomTargetDirFailedMsg, sbomName)
	}

	sourcesbomfilepath := filepath.Join(sbomTmpDir, sbomTmpName)
	targetsbomfilepath := filepath.Join(sbomPath, sbomName)

	err = os.Rename(sourcesbomfilepath, targetsbomfilepath)
	if err != nil {
		return errors.Wrapf(err, mvSBomToTargetDirFailedMsg, sourcesbomfilepath, targetsbomfilepath)
	}
	return nil
}

// mergeSBomFiles - merge sbom files of modules under sbom tmp dir
func mergeSBomFiles(loc *dir.Loc, sbomTmpDir string, sbomName string, sbomSuffix string) (string, error) {
	curtime := time.Now().Format("20230328150313")

	var sbomTmpName string
	if strings.HasSuffix(sbomName, sbom_xml_suffix) {
		sbomTmpName = strings.TrimSuffix(sbomName, xml_suffix) + "_" + curtime + sbom_xml_suffix
	} else if strings.HasSuffix(sbomName, sbom_json_suffix) {
		sbomTmpName = strings.TrimSuffix(sbomName, json_suffix) + "_" + curtime + sbom_xml_suffix
	} else {
		sbomTmpName = sbomName + "_" + curtime + sbom_xml_suffix
	}

	// Get sbom file generate command
	sbomMergeCmds, err := commands.GetSBomsMergeCommand(loc, cyclonedx_cli, sbomTmpDir, sbomTmpName, sbomSuffix)
	if err != nil {
		return "", err
	}

	// exec sbom merge command
	err = executeSBomCommand(sbomMergeCmds)
	if err != nil {
		return "", err
	}
	return sbomTmpName, nil
}

// parseSBomFilePath - parse sbom file path parameter, if it is a relative path, join source path
func parseSBomFilePath(source string, sbomFilePath string) (string, string, string, string) {
	var sbomPath, sbomName, sbomType, sbomSuffix string

	if filepath.IsAbs(sbomFilePath) {
		sbomPath, sbomName = filepath.Split(sbomFilePath)
	} else {
		sbomPath, sbomName = filepath.Split(filepath.Join(source, sbomFilePath))
	}

	sbomType = xml_type
	sbomSuffix = sbom_xml_suffix
	if strings.HasSuffix(sbomName, xml_suffix) {
		sbomType = xml_type
		sbomSuffix = sbom_xml_suffix
	} else if strings.HasSuffix(sbomName, json_suffix) {
		sbomType = json_type
		sbomSuffix = sbom_json_suffix
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

func generateSBomFiles(loc *dir.Loc, mtaObj *mta.MTA, sBomFileTmpDir string, sbomType string, sbomSuffix string) error {
	// (1) sort module by dependency orders
	sortedModuleNames, err := buildops.GetModulesNames(mtaObj)
	if err != nil {
		return err
	}

	// (2) loop modules to generate sbom files
	curtime := time.Now().Format("20230328150313")
	for _, moduleName := range sortedModuleNames {
		module, err := mtaObj.GetModuleByName(moduleName)
		if err != nil {
			return err
		}

		// get sbom file name
		sbomFileName := moduleName + "_" + curtime
		sbomFileFullName := sbomFileName + sbomSuffix

		// get sbom file generate command
		sbomGenCmds, err := commands.GetModuleSBomGenCommands(loc, module, sbomFileName, sbomType, sbomSuffix)
		if err != nil {
			return err
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
