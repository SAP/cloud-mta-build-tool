package artifacts

import (
	"fmt"
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

// ExecuteSBomGenerate - Execute MTA project SBOM generation
func ExecuteSBomGenerate(source string, sbomFilePath string, wdGetter func() (string, error)) error {
	// if sbomFilePath is empty, do not need to generate sbom, return directly
	if strings.TrimSpace(sbomFilePath) == "" {
		return nil
	}

	// (1) validate and parse sbomFilePath
	err := validateSBomFilePath(sbomFilePath)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}
	sbomPath, sbomName, sbomType, sbomSuffix := parseSBomFilePath(sbomFilePath)

	// (2) parse mta.yaml and get mta object
	loc, err := dir.Location(source, "", dir.Dev, []string{}, wdGetter)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	mtaObj, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	// (3) clean and create module sbom generate tmp dir
	sbomTmpDir := loc.GetSBomFileTmpDir(mtaObj)
	err = dir.RemoveIfExist(sbomTmpDir)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}
	err = dir.CreateDirIfNotExist(sbomTmpDir)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	// (4) generation sbom for modules
	err = generateSBomFiles(loc, mtaObj, sbomTmpDir, sbomType, sbomSuffix)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	// (5) merge sbom files
	err = mergeSBomFilesCommand(loc, sbomTmpDir, sbomPath, sbomName, sbomSuffix)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	// (6) generate sbom target dir, mv merged sbom file to target dir
	err = moveSBomToTarget(loc, sbomPath, sbomName, sbomTmpDir)
	if err != nil {
		return errors.Wrapf(err, genSBomFileFailedMsg)
	}

	// (7) clean sbom tmp dir
	// err = dir.RemoveIfExist(sbomTmpDir)
	// if err != nil {
	// 	return errors.Wrapf(err, genSBomFileFailedMsg)
	// }

	return nil
}

func moveSBomToTarget(loc *dir.Loc, sbomPath string, sbomName string, sbomTmpDir string) error {
	sbomTargetPath := filepath.Join(loc.GetSource(), sbomPath)
	err := dir.CreateDirIfNotExist(sbomTargetPath)
	if err != nil {
		return errors.Wrapf(err, createSBomTargetDirFailedMsg, sbomTargetPath)
	}

	sourcesbomfilepath := filepath.Join(sbomTmpDir, sbomName)
	targetsbomfilepath := filepath.Join(sbomTargetPath, sbomName)
	err = os.Rename(sourcesbomfilepath, targetsbomfilepath)
	if err != nil {
		return errors.Wrapf(err, mvSBomToTargetDirFailedMsg, sourcesbomfilepath, targetsbomfilepath)
	}
	return nil
}

func mergeSBomFilesCommand(loc *dir.Loc, sbomTmpDir string, sbomPath string, sbomName string, sbomSuffix string) error {
	// Get sbom file generate command
	sbomMergeCmds, err := commands.GetSBomsMergeCommand(loc, cyclonedx_cli, sbomTmpDir, sbomPath, sbomName, sbomSuffix)
	if err != nil {
		return err
	}

	// exec sbom merge command
	err = executeSBomCommand(sbomMergeCmds)
	if err != nil {
		return err
	}
	return nil
}

func validateSBomFilePath(sbomFilePath string) error {
	if filepath.IsAbs(sbomFilePath) {
		return fmt.Errorf(invalidateSBomFilePath)
	}
	return nil
}

func parseSBomFilePath(sbomFilePath string) (string, string, string, string) {
	filepath, filename := filepath.Split(sbomFilePath)

	filetype := xml_type
	fileSuffix := sbom_xml_suffix
	if strings.HasSuffix(filename, xml_suffix) {
		filetype = xml_type
		fileSuffix = sbom_xml_suffix
	}
	if strings.HasSuffix(filename, json_suffix) {
		filetype = json_type
		fileSuffix = sbom_json_suffix
	}

	return filepath, filename, filetype, fileSuffix
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
		sbomFilefoundPath, err := dir.FindFile(modulePath, sbomFileFullName)
		if err != nil {
			return err
		}
		sbomFileTargetPath := filepath.Join(sBomFileTmpDir, sbomFileFullName)
		err = os.Rename(sbomFilefoundPath, sbomFileTargetPath)
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
