package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
)

var cleanupCmdSrc string
var cleanupCmdTrg string
var cleanupCmdDesc string

var validateCmdSrc string
var validateCmdDesc string
var validateCmdExtensions []string
var validateCmdMode string
var validateCmdStrict string
var validateCmdExclude string

// init - init commands tree and first level commands flags
func init() {

	// Add command to the root
	rootCmd.AddCommand(initCmd, buildCmd, validateCmd, cleanupCmd, provideCmd, generateCmd, moduleCmd, assembleCommand,
		projectCmd, mergeCmd, executeCommand, copyCmd, mtadGenCmd, soloBuildModuleCmd, projectSBomGenCommand, moduleSBomGenCommand)
	// Build module
	provideCmd.AddCommand(provideModuleCmd)
	// generate immutable commands
	generateCmd.AddCommand(metaCmd, mtarCmd)
	// module commands
	moduleCmd.AddCommand(buildModuleCmd, packModuleCmd)
	// project commands
	projectCmd.AddCommand(projectBuildCmd)

	// set flags of cleanup command
	rootCmd.Flags().BoolP("version", "v", false, "Displays the Cloud MTA Build Tool version")
	rootCmd.SetVersionTemplate(rootCmd.Version)
	rootCmd.Flags().BoolP("help", "h", false, "Displays detailed information about the Cloud MTA Build Tool commands; for more information see https://sap.github.io/cloud-mta-build-tool/usage/")

	// set flags of cleanup command
	cleanupCmd.Flags().StringVarP(&cleanupCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is set as default")
	cleanupCmd.Flags().StringVarP(&cleanupCmdTrg, "target", "t", "",
		"The path to the folder in which the temporary artifacts were created; the current path is set as default")
	cleanupCmd.Flags().StringVarP(&cleanupCmdDesc, "desc", "d", "",
		`The MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	cleanupCmd.Flags().BoolP("help", "h", false, `Displays detailed information about the "cleanup" command`)

	// set flags of validation command
	validateCmd.Flags().StringVarP(&validateCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is set as default")
	validateCmd.Flags().StringVarP(&validateCmdMode, "mode", "m", "",
		`The validation mode; supported values: "schema", "semantic" (default)`)
	validateCmd.Flags().StringVarP(&validateCmdDesc, "desc", "d", "",
		`The MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	validateCmd.Flags().StringSliceVarP(&validateCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	validateCmd.Flags().StringVarP(&validateCmdStrict, "strict", "r", "true",
		`If set to true, duplicated fields and fields not defined in the "mta.yaml" schema are reported as errors; if set to false, they are reported as warnings`)
	validateCmd.Flags().StringVarP(&validateCmdExclude, "exclude", "x", "",
		`List of excluded semantic validations; supported validations: "paths", "names", "requires"`)
	validateCmd.Flags().BoolP("help", "h", false, `Displays detailed information about the "validate" command`)

	generateCmd.Flags().BoolP("help", "h", false, `Displays detailed information about the "gen" command`)

}

// generateCmd - Parent of all generation commands
var generateCmd = &cobra.Command{
	Use:    "gen",
	Short:  "Generation commands",
	Long:   "Generation commands",
	Run:    nil,
	Hidden: true,
}

// Parent command - MTA info provider
var provideCmd = &cobra.Command{
	Use:    "provide",
	Short:  "MBT data provider",
	Long:   "MBT data provider",
	Hidden: true,
	Run:    nil,
}

// moduleCmd - Parent of all module commands
var moduleCmd = &cobra.Command{
	Use:    "module",
	Short:  "MBT module commands",
	Long:   "MBT module commands",
	Hidden: true,
	Run:    nil,
}

// moduleCmd - Parent of all module commands
var projectCmd = &cobra.Command{
	Use:    "project",
	Short:  "MBT project commands",
	Long:   "MBT project commands",
	Hidden: true,
	Run:    nil,
}

// Cleanup temp artifacts
var cleanupCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean MBT artifacts",
	Long:  "Clean MBT temporary created artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Remove temp folder
		err := artifacts.ExecuteCleanup(cleanupCmdSrc, cleanupCmdTrg, cleanupCmdDesc, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Validate mta.yaml
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "MBT validation",
	Long:  "MBT validation",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteValidation(validateCmdSrc, validateCmdDesc, validateCmdExtensions, validateCmdMode, validateCmdStrict, validateCmdExclude, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// logError - log errors if any
func logError(err error) {
	if err != nil {
		logs.Logger.Error(err)
	}
}

func cliVersion() string {
	v, _ := version.GetVersionMessage()
	return fmt.Sprintln(v)
}
