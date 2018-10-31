package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/exec"
	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/platform"
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
)

const (
	pathSep    = string(os.PathSeparator)
	dataZip    = pathSep + "data.zip"
	mtarSuffix = ".mtar"
)

var pSourceFlag string
var pTargetFlag string
var pPackModuleFlag string
var pBuildModuleName string
var pValidationFlag string

func init() {
	//set source and target path flags of commands
	setEndpointsFlags(*genMtadCmd, *genMetaCmd, *genMtarCmd, *packCmd, *bModuleCmd, *cleanupCmd)

	// set module flags of module related commands
	packCmd.Flags().StringVarP(&pPackModuleFlag, "module", "m", "", "Provide Module name ")
	bModuleCmd.Flags().StringVarP(&pBuildModuleName, "module", "m", "", "Provide Module name ")

	// set flags of validation command
	validateCmd.Flags().StringVarP(&pValidationFlag, "mode", "m", "", "Provide Validation mode ")
	validateCmd.Flags().StringVarP(&pSourceFlag, "source", "s", "", "Provide MTA source  ")
}

func setEndpointsFlags(commands ...cobra.Command) {
	for _, cmd := range commands {
		cmd.Flags().StringVarP(&pSourceFlag, "source", "s", "", "Provide MTA source ")
		cmd.Flags().StringVarP(&pTargetFlag, "target", "t", "", "Provide MTA target ")
	}
}

// Build module
var bModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Build module",
	Long:  "Build specific module according to the module name",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := buildModule(GetEndPoints(), pBuildModuleName)
		LogError(err)
	},
}

// zip specific module and put the artifacts on the temp folder according
// to the mtar structure, i.e each module have new entry as folder in the mtar folder
// Note - even if the path of the module was changed in the mta.yaml in the mtar the
// the module folder will get the module name
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack module artifacts",
	Long:  "pack the module artifacts after the build process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ep := GetEndPoints()
		modulePath, _, err := getModuleRelativePathAndCommands(ep, pPackModuleFlag)
		if err == nil {
			err = packModule(ep, modulePath, pPackModuleFlag)
		}
		LogError(err)
	},
}

// Generate metadata info from deployment
var genMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateMeta(GetEndPoints())
		LogError(err)
	},
}

// Generate mtar from build artifacts
var genMtarCmd = &cobra.Command{
	Use:   "mtar",
	Short: "generate MTAR",
	Long:  "generate MTAR from the project build artifacts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateMtar(GetEndPoints())
		LogError(err)
	},
}

// Provide mtad.yaml from mta.yaml
var genMtadCmd = &cobra.Command{
	Use:   "mtad",
	Short: "Provide mtad",
	Long:  "Provide deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ep := GetEndPoints()
		mtaStr, err := mta.ReadMta(ep)
		if err == nil {
			err = mta.GenMtad(*mtaStr, ep, func(mtaStr mta.MTA) {
				convertTypes(mtaStr)
			})
		}
		LogError(err)
	},
}

// Validate mta.yaml
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "MBT validation",
	Long:  "MBT validation process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		validateSchema, validateProject, err := getValidationMode(pValidationFlag)
		if err == nil {
			err = validateMtaYaml(GetEndPoints(), validateSchema, validateProject)
		}
		LogError(err)
	},
}

// Cleanup temp artifacts
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove process artifacts",
	Long:  "Remove MTA build process artifacts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting Cleanup process")
		// Remove temp folder
		ep := GetEndPoints()
		err := os.RemoveAll(ep.GetTargetTmpDir())
		if err != nil {
			logs.Logger.Error(err)
		} else {
			logs.Logger.Info("Done")
		}
	},
}

func GetEndPoints() fs.EndPoints {
	return fs.EndPoints{SourcePath: pSourceFlag, TargetPath: pTargetFlag}
}

// generate build metadata artifacts
func generateMeta(ep fs.EndPoints) error {
	return processMta("Metadata creation", ep, []string{}, func(file []byte, args []string) error {
		// Parse MTA file
		m, err := mta.ParseToMta(file)
		if err == nil {
			// Generate meta info dir with required content
			err = mta.GenMetaInfo(ep, *m, args, func(mtaStr mta.MTA) {
				err = convertTypes(mtaStr)
			})
		}
		return err
	})
}

// generate mtar archive from the build artifacts
func generateMtar(ep fs.EndPoints) error {
	logs.Logger.Info(fmt.Sprintf("Generate MTAR. Project path %v, MTAR path %v", ep.GetSource(), ep.GetTarget()))
	return processMta("MTAR generation", ep, []string{}, func(file []byte, args []string) error {
		// read MTA
		m, err := mta.ParseToMta(file)
		if err != nil {
			return err
		}
		// archive building artifacts to mtar
		err = fs.Archive(ep.GetTargetTmpDir(), filepath.Join(ep.GetTarget(), m.Id+mtarSuffix))
		return err
	})
}

// convert types to appropriate target platform types
func convertTypes(mtaStr mta.MTA) error {
	// Load platform configuration file
	platformCfg, err := platform.Parse(platform.PlatformConfig)
	if err == nil {
		// Modify MTAD object according to platform types
		// Todo platform should provided as command parameter
		platform.ConvertTypes(mtaStr, platformCfg, "cf")
	}
	return err
}

// process mta.yaml file
func processMta(processName string, ep fs.EndPoints, args []string, process func(file []byte, args []string) error) error {
	logs.Logger.Info("Starting " + processName)
	mf, err := mta.ReadMtaContent(ep)
	if err == nil {
		err = process(mf, args)
		if err == nil {
			logs.Logger.Info(processName + " finish successfully ")
		}
	} else {
		err = errors.New(fmt.Sprintf("MTA file not found: %s", err))
	}
	return err
}

// pack build module artifacts
func packModule(ep fs.EndPoints, modulePath, moduleName string) error {

	logs.Logger.Info("Pack Module Starts")
	// Get module relative path
	moduleZipPath := ep.GetTargetModuleDir(moduleName)
	logs.Logger.Info(fmt.Sprintf("Module %v will be packed and saved in folder %v", moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return err
	}
	// zipping the build artifacts
	logs.Logger.Infof("Starting execute zipping module %v ", moduleName)
	moduleZipFullPath := moduleZipPath + dataZip
	if err = fs.Archive(ep.GetSourceModuleDir(modulePath), moduleZipFullPath); err != nil {
		err = errors.New(fmt.Sprintf("Error occurred during ZIP module %v creation, error: %s  ", moduleName, err))
		removeErr := os.RemoveAll(ep.GetTargetTmpDir())
		if removeErr != nil {
			err = errors.New(fmt.Sprintf("Error occured during directory %s removal failed %s. %s", ep.GetTargetTmpDir(), err, removeErr))
		}
	} else {
		logs.Logger.Infof("Execute zipping module %v finished successfully ", moduleName)
	}
	return err
}

// convert validation mode flag to validation process flags
func getValidationMode(validationFlag string) (bool, bool, error) {
	switch validationFlag {
	case "":
		return true, true, nil
	case "schema":
		return true, false, nil
	case "project":
		return false, true, nil
	}
	return false, false, errors.New("wrong argument of validation mode. Expected one of [all, schema, project]")
}

// Validate MTA yaml
func validateMtaYaml(ep fs.EndPoints, validateSchema bool, validateProject bool) error {
	if validateProject || validateSchema {
		logs.Logger.Info("Starting MTA Yaml validation")

		// Read MTA yaml content
		yamlContent, err := mta.ReadMtaContent(ep)

		if err != nil {
			return errors.New("MTA validation failed. " + err.Error())
		}
		projectPath := ep.GetSource()

		// validate mta content
		issues := mta.Validate(yamlContent, projectPath, validateSchema, validateProject)
		if len(issues) == 0 {
			logs.Logger.Info("MTA Yaml is valid")
		} else {
			return errors.New("MTA Yaml is  invalid. Issues: \n" + issues.String())
		}
	}
	return nil
}

// Get module relative path from mta.yaml and
// commands (with resolved paths) configured for the module type
func getModuleRelativePathAndCommands(ep fs.EndPoints, module string) (string, []string, error) {
	mtaObj, err := mta.ReadMta(ep)
	if err != nil {
		return "", nil, err
	}
	// Get module respective command's to execute
	return moduleCmd(*mtaObj, module)
}

func buildModule(ep fs.EndPoints, module string) error {

	logs.Logger.Info("Start building module: ", module)
	// Get module respective command's to execute
	moduleRelPath, mCmd, err := getModuleRelativePathAndCommands(ep, module)
	if err != nil {
		return err
	}
	modulePath := ep.GetSourceModuleDir(moduleRelPath)

	// Get module commands
	commands := cmdConverter(modulePath, mCmd)

	// Execute child-process with module respective commands
	err = exec.Execute(commands)
	if err != nil {
		return err
	}
	// Pack the modules build artifacts (include node modules)
	// into the artifactsPath dir as data zip
	return packModule(ep, moduleRelPath, module)
}

// Get commands for specific module type
func moduleCmd(mta mta.MTA, moduleName string) (string, []string, error) {
	for _, m := range mta.Modules {
		if m.Name == moduleName {
			commandProvider, err := builders.CommandProvider(*m)
			if err != nil {
				return "", nil, err
			}
			return m.Path, commandProvider.Command, nil
		}
	}
	return "", nil, errors.New(fmt.Sprintf("Module %v not defined in MTA", moduleName))
}

// path and commands to execute
func cmdConverter(mPath string, cmdList []string) [][]string {
	var cmd [][]string
	for i := 0; i < len(cmdList); i++ {
		cmd = append(cmd, append([]string{mPath}, strings.Split(cmdList[i], " ")...))
	}
	return cmd
}
