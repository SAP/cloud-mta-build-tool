package commands

import (
	"os"

	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
	"github.com/spf13/cobra"
)

var buildTargetFlag string
var validationFlag string

// Parent command
var build = &cobra.Command{
	Use:   "build",
	Short: "Build Project",
	Long:  "Build MTA project",
	Run:   nil,
}

// Parent command
var execute = &cobra.Command{
	Use:   "execute",
	Short: "Execute step",
	Long:  "Execute standalone step as part of the build process",
	Run:   nil,
}

// Parent command
var provides = &cobra.Command{
	Use:   "provide",
	Short: "MBT data provider",
	Long:  "MBT data provider",
	Run:   nil,
}

// Parent command
var validate = &cobra.Command{
	Use:   "validate",
	Short: "MBT validation",
	Long:  "MBT validation process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		validateSchema, validateProject := getValidationMode(validationFlag)
		validateMtaYaml("", "mta.yaml", validateSchema, validateProject)
	},
}

func getValidationMode(validationFlag string) (bool, bool) {
	switch true {
	case validationFlag == "":
		return true, true
	case validationFlag == "schema":
		return true, false
	case validationFlag == "project":
		return false, true
	}
	logs.Logger.Error("Wrong argument of validation mode. Expected one of [all, schema, project")
	return false, false
}

func validateMtaYaml(yamlPath string, yamlFilename string, validateSchema bool, validateProject bool) {
	if !validateProject && !validateSchema {
		return
	}
	logs.Logger.Info("Starting MTA Yaml validation")
	yamlContent, err := mta.ReadMtaContent(yamlPath, yamlFilename)
	var wd string
	if err == nil {
		wd, err = os.Getwd()
	}
	if err != nil {
		logs.Logger.Error("MTA validation failed. " + err.Error())
	} else {
		issues := mta.Validate(yamlContent, wd, validateSchema, validateProject)
		valid := len(issues) == 0
		if valid {
			logs.Logger.Info("MTA Yaml is valid")
		} else {
			logs.Logger.Info("MTA Yaml is  invalid. Issues: \n", issues)
		}
	}
}

func init() {

	// Build module
	provides.AddCommand(pModule)
	// Provide module
	build.AddCommand(bModule)
	// execute immutable commands
	execute.AddCommand(prepare, pack, genMeta, pMtad, genMtar, cleanup, validate)
	// Add command to the root
	rootCmd.AddCommand(provides, build, execute, initProcess)
	// build target flags
	build.Flags().StringVarP(&buildTargetFlag, "target", "t", "", "Build for specified environment ")
	// validation flags , can be used for multiple scenario
	validate.Flags().StringVarP(&validationFlag, "mode", "m", "", "Validation mode ")
}
