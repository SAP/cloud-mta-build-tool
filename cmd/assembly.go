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
var assembleCmdMtarName string
var assembleCmdParallel string

// Assemble the MTA project post-build artifacts, without any build process
var assemblyCommand = &cobra.Command{
	Use:       "assemble",
	Short:     "Assembles MTA Archive",
	Long:      "Assembles MTA Archive",
	ValidArgs: []string{"Deployment descriptor location"},
	Args:      cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.Assembly(assembleCmdSrc, assembleCmdTrg, defaultPlatform, assembleCmdMtarName, assembleCmdParallel, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	assemblyCommand.Flags().StringVarP(&assembleCmdSrc,
		"source", "s", "", "the path to the MTA project; the current path is set as the default")
	assemblyCommand.Flags().StringVarP(&assembleCmdTrg,
		"target", "t", "", "the path to the MBT results folder; the current path is set as the default")
	assemblyCommand.Flags().StringVarP(&assembleCmdMtarName,
		"mtar", "m", "", "the archive name")
	assemblyCommand.Flags().StringVarP(&assembleCmdParallel,
		"parallel", "p", "false", "if true content copying will run in parallel")

}
