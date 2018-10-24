package commands

import (
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/tpl"
)

var initMode string

func init() {
	initProcessCmd.Flags().StringVarP(&initMode, "mode", "m", "", "Mode of Makefile generation - default/verbose")
}

var initProcessCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate Makefile",
	Long:  "Generate Makefile as manifest which describe's the build process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.Make(initMode)
		LogError(err)
	},
}
