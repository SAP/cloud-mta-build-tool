package commands

import (
	"github.com/spf13/cobra"
)

// Parent commands
var build = &cobra.Command{
	Use:   "build",
	Short: "Build Project",
	Long:  "Build MTA project",
	Run:   nil,
}

// Execute small building blocks
var execute = &cobra.Command{
	Use:   "execute",
	Short: "Execute step",
	Long:  "Execute standalone step as part of the build process",
	Run:   nil,
}

func init() {
	build.AddCommand(cfBuild, neoBuild, html5)
	execute.AddCommand(prepare, copyModule, pack, genMeta, genMtar, cleanup)
	rootCmd.AddCommand(build, execute, initProcess)
}
