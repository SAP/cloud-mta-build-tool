package commands

import (
	"github.com/spf13/cobra"
)

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
	provides.AddCommand(pm)
	build.AddCommand(cfBuild, neoBuild)
	execute.AddCommand(prepare, copyModule, pack, genMeta, genMtar, cleanup)
	rootCmd.AddCommand(provides, build, execute, initProcess)
}
