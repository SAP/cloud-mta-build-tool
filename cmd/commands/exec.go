package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

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

var sourceMtadFlag string
var targetMtadFlag string
var sourceMetaFlag string
var targetMetaFlag string
var sourceMtarFlag string
var targetMtarFlag string
var sourcePackFlag string
var targetPackFlag string
var sourceBModuleFlag string
var targetBModuleFlag string
var sourceCleanupFlag string
var targetCleanupFlag string
var sourceValidateFlag string

var pPackModuleFlag string
var pBuildModuleNameFlag string
var pValidationFlag string

var descriptorMtadFlag string
var descriptorMtarFlag string
var descriptorMetaFlag string
var descriptorPackFlag string
var descriptorBModuleFlag string
var descriptorCleanupFlag string
var descriptorValidateFlag string

func init() {

	// set source and target path flags of commands
	genMtadCmd.Flags().StringVarP(&sourceMtadFlag, "source", "s", "", "Provide MTA source ")
	genMtadCmd.Flags().StringVarP(&targetMtadFlag, "target", "t", "", "Provide MTA target ")
	genMetaCmd.Flags().StringVarP(&sourceMetaFlag, "source", "s", "", "Provide MTA source ")
	genMetaCmd.Flags().StringVarP(&targetMetaFlag, "target", "t", "", "Provide MTA target ")
	genMtarCmd.Flags().StringVarP(&sourceMtarFlag, "source", "s", "", "Provide MTA source ")
	genMtarCmd.Flags().StringVarP(&targetMtarFlag, "target", "t", "", "Provide MTA target ")
	packCmd.Flags().StringVarP(&sourcePackFlag, "source", "s", "", "Provide MTA source ")
	packCmd.Flags().StringVarP(&targetPackFlag, "target", "t", "", "Provide MTA target ")
	bModuleCmd.Flags().StringVarP(&sourceBModuleFlag, "source", "s", "", "Provide MTA source ")
	bModuleCmd.Flags().StringVarP(&targetBModuleFlag, "target", "t", "", "Provide MTA target ")
	cleanupCmd.Flags().StringVarP(&sourceCleanupFlag, "source", "s", "", "Provide MTA source ")
	cleanupCmd.Flags().StringVarP(&targetCleanupFlag, "target", "t", "", "Provide MTA target ")
	validateCmd.Flags().StringVarP(&sourceValidateFlag, "source", "s", "", "Provide MTA source  ")

	// set module flags of module related commands
	packCmd.Flags().StringVarP(&pPackModuleFlag, "module", "m", "", "Provide Module name ")
	bModuleCmd.Flags().StringVarP(&pBuildModuleNameFlag, "module", "m", "", "Provide Module name ")

	// set flags of validation command
	validateCmd.Flags().StringVarP(&pValidationFlag, "mode", "m", "", "Provide Validation mode ")

	// set mta descriptor flag of commands
	genMtadCmd.Flags().StringVarP(&descriptorMtadFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
	genMetaCmd.Flags().StringVarP(&descriptorMetaFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
	genMtarCmd.Flags().StringVarP(&descriptorMtarFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
	packCmd.Flags().StringVarP(&descriptorPackFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
	bModuleCmd.Flags().StringVarP(&descriptorBModuleFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
	cleanupCmd.Flags().StringVarP(&descriptorCleanupFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
	validateCmd.Flags().StringVarP(&descriptorValidateFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
}

// Build module
var bModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Build module",
	Long:  "Build specific module according to the module name",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := fs.ValidateDeploymentDescriptor(descriptorBModuleFlag)
		if err == nil {
			ep := GetLocationParameters(sourceBModuleFlag, targetBModuleFlag, descriptorBModuleFlag)
			err = buildModule(&ep, pBuildModuleNameFlag)
		}
		return err
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
		ep := GetLocationParameters(sourcePackFlag, targetPackFlag, descriptorPackFlag)
		modulePath, _, err := getModuleRelativePathAndCommands(&ep, pPackModuleFlag)
		if err == nil {
			err = packModule(&ep, modulePath, pPackModuleFlag)
		}
		logError(err)
	},
}

// Generate metadata info from deployment
var genMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := fs.ValidateDeploymentDescriptor(descriptorMetaFlag)
		if err == nil {
			ep := GetLocationParameters(sourceMetaFlag, targetMetaFlag, descriptorMetaFlag)
			err = generateMeta(&ep)
		}
		logErrorExt(err, "META generation failed")
	},
}

// Generate mtar from build artifacts
var genMtarCmd = &cobra.Command{
	Use:   "mtar",
	Short: "generate MTAR",
	Long:  "generate MTAR from the project build artifacts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := fs.ValidateDeploymentDescriptor(descriptorMtarFlag)
		if err == nil {
			ep := GetLocationParameters(sourceMtarFlag, targetMtarFlag, descriptorMtarFlag)
			err = generateMtar(&ep)
		}
		logErrorExt(err, "MTAR generation failed")
	},
}

// Provide mtad.yaml from mta.yaml
var genMtadCmd = &cobra.Command{
	Use:   "mtad",
	Short: "Provide mtad",
	Long:  "Provide deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := fs.ValidateDeploymentDescriptor(descriptorMtadFlag)
		if err != nil {
			logErrorExt(err, "MTAD generation failed")
			return err
		}
		ep := GetLocationParameters(sourceMtadFlag, targetMtadFlag, descriptorMtadFlag)
		// TODO if descriptor == "dep" -> Copy mtad
		mtaStr, err := mta.ReadMta(&ep)
		if err == nil {
			err = mta.GenMtad(mtaStr, &ep, func(mtaStr *mta.MTA) {
				e := convertTypes(*mtaStr);if e != nil {
					logErrorExt(err, "MTAD generation failed")
				}
			})
		}
		logErrorExt(err, "MTAD generation failed")
		return err
	},
}

// Validate mta.yaml
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "MBT validation",
	Long:  "MBT validation process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := fs.ValidateDeploymentDescriptor(descriptorValidateFlag)
		if err != nil {
			logErrorExt(err, "MBT Validation failed")
			return
		}
		validateSchema, validateProject, err := getValidationMode(pValidationFlag)
		if err == nil {
			ep := GetLocationParameters(sourceValidateFlag, sourceValidateFlag, descriptorValidateFlag)
			err = validateMtaYaml(&ep, validateSchema, validateProject)
		}
		logErrorExt(err, "MBT Validation failed")
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
		ep := GetLocationParameters(sourceCleanupFlag, targetCleanupFlag, descriptorCleanupFlag)
		err := os.RemoveAll(ep.GetTargetTmpDir())
		if err != nil {
			logs.Logger.Error(err)
		} else {
			logs.Logger.Info("Done")
		}
	},
}

// GetLocationParameters - provides location parameters of MTA
func GetLocationParameters(sourceFlag, targetFlag, descriptor string) fs.MtaLocationParameters {
	var mtaFilename string
	if descriptor == "dev" || descriptor == "" {
		mtaFilename = "mta.yaml"
		descriptor = "dev"
	} else {
		mtaFilename = "mtad.yaml"
		descriptor = "dep"
	}
	return fs.MtaLocationParameters{SourcePath: sourceFlag, TargetPath: targetFlag, MtaFilename: mtaFilename, Descriptor: descriptor}
}

// generate build metadata artifacts
func generateMeta(ep *fs.MtaLocationParameters) error {
	return processMta("Metadata creation", ep, []string{}, func(file []byte, args []string) error {
		// Parse MTA file
		m, err := mta.ParseToMta(file)
		if err == nil {
			// Generate meta info dir with required content
			err = mta.GenMetaInfo(ep, m, args, func(mtaStr *mta.MTA) {
				err = convertTypes(*mtaStr)
			})
		}
		return err
	})
}

// generate mtar archive from the build artifacts
func generateMtar(ep *fs.MtaLocationParameters) error {
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
func processMta(processName string, ep *fs.MtaLocationParameters, args []string, process func(file []byte, args []string) error) error {
	logs.Logger.Info("Starting " + processName)
	mf, err := mta.ReadMtaContent(ep)
	if err == nil {
		err = process(mf, args)
		if err == nil {
			logs.Logger.Info(processName + " finish successfully ")
		}
	} else {
		err = errors.Wrap(err, "MTA file not found")
	}
	return err
}

// pack build module artifacts
func packModule(ep *fs.MtaLocationParameters, modulePath, moduleName string) error {

	logs.Logger.Info("Pack Module Starts")
	// Get module relative path
	moduleZipPath := ep.GetTargetModuleDir(moduleName)
	logs.Logger.Info(fmt.Sprintf("Module %v will be packed and saved in folder %v", moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Error occurred during creation of directory of ZIP module %v", moduleName)
	}
	// zipping the build artifacts
	logs.Logger.Infof("Starting execute zipping module %v ", moduleName)
	moduleZipFullPath := moduleZipPath + dataZip
	if err = fs.Archive(ep.GetSourceModuleDir(modulePath), moduleZipFullPath); err != nil {
		err = errors.Wrapf(err, "Error occurred during ZIP module %v creation", moduleName)
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
func validateMtaYaml(ep *fs.MtaLocationParameters, validateSchema bool, validateProject bool) error {
	if validateProject || validateSchema {
		logs.Logger.Infof("Starting %v validation", ep.MtaFilename)

		// Read MTA yaml content
		yamlContent, err := mta.ReadMtaContent(ep)

		if err != nil {
			return errors.Wrap(err, "MTA validation failed")
		}
		projectPath := ep.GetSource()

		// validate mta content
		issues := mta.Validate(yamlContent, projectPath, validateSchema, validateProject)
		if len(issues) == 0 {
			logs.Logger.Infof("%v is valid", ep.MtaFilename)
		} else {
			return errors.New(fmt.Sprintf("%v is invalid. Issues: \n%v", ep.MtaFilename, issues.String()))
		}
	}

	return nil
}

// Get module relative path from mta.yaml and
// commands (with resolved paths) configured for the module type
func getModuleRelativePathAndCommands(ep *fs.MtaLocationParameters, module string) (string, []string, error) {
	mtaObj, err := mta.ReadMta(ep)
	if err != nil {
		return "", nil, err
	}
	// Get module respective command's to execute
	return moduleCmd(mtaObj, module)
}

func buildModule(ep *fs.MtaLocationParameters, module string) error {

	if ep.IsDeploymentDescriptor() {
		return nil
	}

	logs.Logger.Info("Start building module: ", module)

	err := processDependencies(ep, module)
	if err != nil {
		return err
	}

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
func moduleCmd(mta *mta.MTA, moduleName string) (string, []string, error) {
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

func processDependencies(ep *fs.MtaLocationParameters, moduleName string) error {
	mtaObj, err := mta.ReadMta(ep)
	if err != nil {
		return err
	}
	module, err := mtaObj.GetModuleByName(moduleName)
	if err != nil {
		return err
	}
	if module.Requires != nil {
		for _, req := range module.BuildParams.Requires {
			err := req.ProcessRequirements(ep, mtaObj, module.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
