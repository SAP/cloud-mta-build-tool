package commands

import (
	"cloud-mta-build-tool/internal/logs"
	"github.com/spf13/cobra"
)

var buildTargetFlag string
var validationFlag string

func init() {

	// Build module
	provideCmd.AddCommand(pModuleCmd)
	// Provide module
	buildCmd.AddCommand(bModuleCmd)
	// execute immutable commands
	executeCmd.AddCommand(packCmd, genMetaCmd, genMtadCmd, genMtarCmd, cleanupCmd, validateCmd)
	// Add command to the root
	rootCmd.AddCommand(provideCmd, buildCmd, executeCmd, initProcessCmd)
	// build command target flags
	buildCmd.Flags().StringVarP(&buildTargetFlag, "target", "t", "", "Build for specified environment ")
	// validation flags , can be used for multiple scenario
	validateCmd.Flags().StringVarP(&validationFlag, "mode", "m", "", "Validation mode ")
}

// Parent command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Project",
	Long:  "Build MTA project",
	Run:   nil,
}

// Parent command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute step",
	Long:  "Execute standalone step as part of the build process",
	Run:   nil,
}

// Parent command
var provideCmd = &cobra.Command{
	Use:   "provide",
	Short: "MBT data provider",
	Long:  "MBT data provider",
	Run:   nil,
}

// LogError - log errors if any
func LogError(err error) {
	if err != nil {
		logs.Logger.Error(err)
	}
}
