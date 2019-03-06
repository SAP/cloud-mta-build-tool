package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

var projectBuildCmdSrc string
var projectBuildCmdDesc string
var projectBuildCmdPhase string

func init() {
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdSrc,
		"source", "s", "", "the path to the MTA project; the current path is default")
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdDesc,
		"desc", "d", "", "the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdPhase,
		"phase", "p", "", "the project build phase; supported values: pre and post")
}

// projectBuildCmd - Runs the mta project pre and post build processes
var projectBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run the MTA project pre and post build commands",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteProjectBuild(projectBuildCmdSrc, projectBuildCmdDesc, projectBuildCmdPhase, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}
