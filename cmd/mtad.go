package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

// mtad gen - command flags
var mtadGenCmdSrc string
var mtadGenCmdTrg string
var mtadGenCmdExtensions []string
var mtadGenCmdPlatform string

// Provide mtad.yaml from mta.yaml
var mtadGenCmd = &cobra.Command{
	Use:   "mtad-gen",
	Short: "Generates an 'mtad.yaml' file",
	Long:  "Generates a deployment descriptor ('mtad.yaml') file from the 'mta.yaml' file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteMtadGen(mtadGenCmdSrc, mtadGenCmdTrg, mtadGenCmdExtensions, mtadGenCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	// set flags of mtad gen command
	mtadGenCmd.Flags().StringVarP(&mtadGenCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is set as default")
	mtadGenCmd.Flags().StringVarP(&mtadGenCmdTrg, "target", "t",
		"", "The path to the folder in which the 'mtad.yaml' file is generated; the current path is set as default")
	mtadGenCmd.Flags().StringSliceVarP(&mtadGenCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	mtadGenCmd.Flags().StringVarP(&mtadGenCmdPlatform, "platform", "p", "cf",
		`The deployment platform; supported platforms: "cf", "xsa", "neo"`)
	mtadGenCmd.Flags().BoolP("help", "h", false, `Displays detailed information about the 'mtad gen' command`)
}
