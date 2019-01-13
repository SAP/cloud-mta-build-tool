package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/artifacts"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/version"
)

var cleanupCmdSrc string
var cleanupCmdTrg string
var cleanupCmdDesc string

var validateCmdSrc string
var validateCmdDesc string
var validateCmdMode string

// init - init commands tree and first level commands flags
func init() {

	// Add command to the root
	rootCmd.AddCommand(versionCmd, initCmd, validateCmd, cleanupCmd, provideCmd, generateCmd, moduleCmd, assemblyCommand)
	// Build module
	provideCmd.AddCommand(provideModuleCmd)
	// generate immutable commands
	generateCmd.AddCommand(metaCmd, mtadCmd, mtarCmd)
	// module commands
	moduleCmd.AddCommand(buildModuleCmd, packModuleCmd)

	// set flags of cleanup command
	cleanupCmd.Flags().StringVarP(&cleanupCmdSrc, "source", "s", "", "Provide MTA source ")
	cleanupCmd.Flags().StringVarP(&cleanupCmdTrg, "target", "t", "", "Provide MTA target ")
	cleanupCmd.Flags().StringVarP(&cleanupCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")

	// set flags of validation command
	validateCmd.Flags().StringVarP(&validateCmdSrc, "source", "s", "", "Provide MTA source  ")
	validateCmd.Flags().StringVarP(&validateCmdMode, "mode", "m", "", "Provide Validation mode ")
	validateCmd.Flags().StringVarP(&validateCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
}

// generateCmd - Parent of all generation commands
var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generation step",
	Long:  "Execute standalone step as part of the build process",
	Run:   nil,
}

// Parent command - MTA info provider
var provideCmd = &cobra.Command{
	Use:   "provide",
	Short: "MBT data provider",
	Long:  "MBT data provider",
	Run:   nil,
}

// moduleCmd - Parent of all module commands
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "MBT module commands",
	Long:  "MBT data provider",
	Run:   nil,
}

// Parent command - CLI Version provider
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "MBT version",
	Long:  "MBT version",
	Run: func(cmd *cobra.Command, args []string) {
		err := printCliVersion()
		logError(err)
	},
}

// Validate mta.yaml
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "MBT validation",
	Long:  "MBT validation process",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteValidation(validateCmdSrc, validateCmdDesc, validateCmdMode, os.Getwd)
		logError(err)
		return err
	},
	SilenceErrors: false,
	SilenceUsage:  true,
}

// Cleanup temp artifacts
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove process artifacts",
	Long:  "Remove MTA build process artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Remove temp folder
		err := artifacts.ExecuteCleanup(cleanupCmdSrc, cleanupCmdTrg, cleanupCmdDesc, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}

func printCliVersion() error {
	v, err := version.GetVersion()
	if err == nil {
		fmt.Println(v.CliVersion)
	}
	return err
}

// logError - log errors if any
func logError(err error) {
	if err != nil {
		logs.Logger.Error(err)
	}
}
