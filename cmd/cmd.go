package commands

import (
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
var validateCmdMode string
var validateCmdStrict string

// init - init commands tree and first level commands flags
func init() {

	// Add command to the root
	rootCmd.AddCommand(initCmd, validateCmd, cleanupCmd, provideCmd, generateCmd, moduleCmd, assemblyCommand)
	// Build module
	provideCmd.AddCommand(provideModuleCmd)
	// generate immutable commands
	generateCmd.AddCommand(metaCmd, mtadCmd, mtarCmd)
	// module commands
	moduleCmd.AddCommand(buildModuleCmd, packModuleCmd)

	// set flags of cleanup command
	rootCmd.Flags().BoolP("version", "v", false, "version for MBT")

	// set flags of cleanup command
	cleanupCmd.Flags().StringVarP(&cleanupCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is default")
	cleanupCmd.Flags().StringVarP(&cleanupCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is default")
	cleanupCmd.Flags().StringVarP(&cleanupCmdDesc, "desc", "d", "",
		"the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")

	// set flags of validation command
	validateCmd.Flags().StringVarP(&validateCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is default")
	validateCmd.Flags().StringVarP(&validateCmdMode, "mode", "m", "",
		"the validation mode; supported values: schema, semantic (default)")
	validateCmd.Flags().StringVarP(&validateCmdDesc, "desc", "d", "",
		"the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")
	validateCmd.Flags().StringVarP(&validateCmdStrict, "strict", "r", "",
		`if set to true, duplicated fields and fields not defined in the "mta.yaml" schema are reported as errors; if set to false, they are reported as warnings`)
}

// generateCmd - Parent of all generation commands
var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "generation commands",
	Long:  "generation commands",
	Run:   nil,
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

// Cleanup temp artifacts
var cleanupCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean MBT artifacts",
	Long:  "clean MBT temporary created artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Remove temp folder
		err := artifacts.ExecuteCleanup(cleanupCmdSrc, cleanupCmdTrg, cleanupCmdDesc, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	Hidden:        true,
	SilenceErrors: true,
}

// Validate mta.yaml
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "MBT validation",
	Long:  "MBT validation",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteValidation(validateCmdSrc, validateCmdDesc, validateCmdMode, validateCmdStrict, os.Getwd)
		logError(err)
		return err
	},
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
	v, _ := version.GetVersion()
	return v.CliVersion
}
