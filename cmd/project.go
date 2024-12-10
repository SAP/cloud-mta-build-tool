package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

var projectBuildCmdSrc string
var projectBuildCmdMtaYamlFilename string
var projectBuildCmdTrg string
var projectBuildCmdDesc string
var projectBuildCmdExtensions []string
var projectBuildCmdPhase string

func init() {
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdSrc,
		"source", "s", "", "The path to the MTA project; the current path is set as default")
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdMtaYamlFilename,
		"filename", "f", "", "The mta yaml filename of MTA project; the mta.yaml is set as default")
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdTrg,
		"target", "t", "", "The path to the folder in which the temporary artifacts of the project build are created; the current path is set as default")
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdDesc,
		"desc", "d", "", `The MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	projectBuildCmd.Flags().StringSliceVarP(&projectBuildCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	projectBuildCmd.Flags().StringVarP(&projectBuildCmdPhase,
		"phase", "p", "", `The project build phase; supported values: "pre" and "post"`)
}

// projectBuildCmd - Runs the mta project pre and post build processes
var projectBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run the MTA project pre and post build commands",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteProjectBuild(projectBuildCmdSrc, projectBuildCmdMtaYamlFilename, projectBuildCmdTrg, projectBuildCmdDesc, projectBuildCmdExtensions, projectBuildCmdPhase, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}
