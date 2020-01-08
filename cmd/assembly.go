package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

const (
	defaultPlatform string = "cf"
)

var assembleCmdSrc string
var assembleCmdTrg string
var assembleCmdExtensions []string
var assembleCmdMtarName string
var assembleCmdParallel string

// Assemble the MTA project post-build artifacts, without any build process
var assembleCommand = &cobra.Command{
	Use:   "assemble",
	Short: "Generates an MTA archive according to the MTA deployment descriptor (mtad.yaml)",
	Long:  "Generates an MTA archive according to the MTA deployment descriptor (mtad.yaml)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.Assembly(assembleCmdSrc, assembleCmdTrg, assembleCmdExtensions, defaultPlatform, assembleCmdMtarName, assembleCmdParallel, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	assembleCommand.Flags().StringVarP(&assembleCmdSrc,
		"source", "s", "", "The path to the MTA project; the current path is set as default")
	assembleCommand.Flags().StringVarP(&assembleCmdTrg,
		"target", "t", "", `The path to the generated MTAR file; the path to the "mta_archives" subfolder of the current folder is set as default`)
	assembleCommand.Flags().StringSliceVarP(&assembleCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	assembleCommand.Flags().StringVarP(&assembleCmdMtarName,
		"mtar", "m", "", "The archive name")
	assembleCommand.Flags().StringVarP(&assembleCmdParallel,
		"parallel", "p", "true", "If true content copying will run in parallel")
	_ = assembleCommand.Flags().MarkHidden("parallel")
	assembleCommand.Flags().BoolP("help", "h", false, `Displays detailed information about the "assemble" command`)

}
