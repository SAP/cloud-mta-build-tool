package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
)

// flags of init command
var initCmdSrc string
var initCmdTrg string
var initCmdMode string

// init flags of init command
func init() {
	initCmd.Flags().StringVarP(&initCmdSrc, "source", "s", "", "the path to the MTA project; the current path is default")
	initCmd.Flags().StringVarP(&initCmdTrg, "target", "t", "", "the path to the generated Makefile folder; the current path is default")
	initCmd.Flags().StringVarP(&initCmdMode, "mode", "m", "", "the mode of the Makefile generation; supported values: default and verbose")
	initCmd.Flags().MarkHidden("mode")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generates Makefile",
	Long:  "generates Makefile as a manifest file that describes the build process of the MTA project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.ExecuteMake(initCmdSrc, initCmdTrg, initCmdMode, os.Getwd)
		logError(err)
	},
}
