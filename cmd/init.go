package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
)

// flags of init command
var initCmdSrc string
var initCmdTrg string
var initCmdDesc string
var initCmdMode string

// init flags of init command
func init() {
	initCmd.Flags().StringVarP(&initCmdSrc, "source", "s", "", "Provide MTA source")
	initCmd.Flags().StringVarP(&initCmdTrg, "target", "t", "", "Provide MTA target")
	initCmd.Flags().StringVarP(&initCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	initCmd.Flags().StringVarP(&initCmdMode, "mode", "m", "", "Mode of Makefile generation - default/verbose")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generates Makefile",
	Long:  "generates Makefile as a manifest file that describes the build process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.ExecuteMake(initCmdSrc, initCmdTrg, initCmdDesc, initCmdMode, os.Getwd)
		logError(err)
	},
}
