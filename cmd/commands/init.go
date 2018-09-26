package commands

import (
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/tpl"
)

var initProcess = &cobra.Command{
	Use:   "init",
	Short: "Generate Makefile",
	Long:  "Generate Makefile as manifest which describe's the build process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		if err := tpl.Make(initMode); err != nil {
			logs.Logger.Error(err)
		}
	},
}

func init() {
	initProcess.Flags().StringVarP(&initMode, "mode", "m", "", "Mode of Makefile generation - default/verbose")
}
