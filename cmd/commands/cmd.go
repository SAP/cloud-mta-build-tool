package commands

import (
	"github.com/spf13/cobra"
)

var initMode string
var buildTargetEnv string

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

func init() {
	// build target flags
	build.Flags().StringVarP(&buildTargetEnv, "target", "t", "", "Build for specified environment ")
	// Build module
	provides.AddCommand(pModule)
	// Provide module
	build.AddCommand(bModule)
	// execute immutable commands
	execute.AddCommand(prepare, pack, genMeta, genMtar, cleanup)
	// Add command to the root
	rootCmd.AddCommand(provides, build, execute, initProcess)
}
