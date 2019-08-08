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
var assembleExtensions []string
var assembleCmdMtarName string
var assembleCmdParallel string

// Assemble the MTA project post-build artifacts, without any build process
var assembleCommand = &cobra.Command{
	Use:       "assemble",
	Short:     "Assembles MTA Archive",
	Long:      "Assembles MTA Archive",
	ValidArgs: []string{"Deployment descriptor location"},
	Args:      cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.Assembly(assembleCmdSrc, assembleCmdTrg, assembleExtensions, defaultPlatform, assembleCmdMtarName, assembleCmdParallel, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	assembleCommand.Flags().StringVarP(&assembleCmdSrc,
		"source", "s", "", "the path to the MTA project; the current path is set as the default")
	assembleCommand.Flags().StringVarP(&assembleCmdTrg,
		"target", "t", "", "the path to the MBT results folder; the current path is set as the default")
	assembleCommand.Flags().StringSliceVarP(&assembleExtensions, "extensions", "e", nil,
		"the MTA extension descriptors")
	assembleCommand.Flags().StringVarP(&assembleCmdMtarName,
		"mtar", "m", "", "the archive name")
	assembleCommand.Flags().StringVarP(&assembleCmdParallel,
		"parallel", "p", "true", "if true content copying will run in parallel")
	_ = assembleCommand.Flags().MarkHidden("parallel")
	assembleCommand.Flags().BoolP("help", "h", false, `prints detailed information about the "assemble" command`)

}
