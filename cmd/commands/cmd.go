package commands

import (
	"cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
	"github.com/spf13/cobra"
)

var initMode string
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
	Run: func(cmd *cobra.Command, args []string) {
		validateSchema, validateProject := getValidationMode(args)
		validateMtaYaml(dir.GetPath(), "mta.yaml", validateSchema, validateProject)
	},
}

func getValidationMode(args []string) (bool, bool) {
	switch true {
	case len(args) == 0 || args[0] == "all":
		return true, true
	case args[0] == "schema":
		return true, false
	case args[0] == "project":
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
	source := mta.Source{yamlPath, yamlFilename}
	yamlContent, err := source.ReadExtFile()
	if err != nil {
		logs.Logger.Error(err)
	} else {
		issues := mta.Validate(yamlContent, dir.GetPath(), validateSchema, validateProject)
		logs.Logger.Info("MTA Yaml is %t", issues)
	}
}

func init() {

	// Build module
	provides.AddCommand(pModule)
	// Provide module
	build.AddCommand(bModule)
	// execute immutable commands
	execute.AddCommand(prepare, pack, genMeta, genMtar, cleanup)
	// Add command to the root
	rootCmd.AddCommand(provides, build, execute, initProcess)
	// build target flags
	build.Flags().StringVarP(&buildTargetFlag, "target", "t", "", "Build for specified environment ")
	// validation flags , can be used for multiple scenario
	validate.Flags().StringVarP(&validationFlag, "validate", "v", "", "Validation process ")
}
