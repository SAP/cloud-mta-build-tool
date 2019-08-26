package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

var mergeCmdSrc string
var mergeCmdTrg string
var mergeCmdExtensions []string
var mergeCmdName string

// Merge mta.yaml with mta extension files and write the result.
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: `Merges the "mta.yaml" file with the MTA extension descriptors`,
	Long:  `Merges the "mta.yaml" file with the MTA extension descriptors`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteMerge(mergeCmdSrc, mergeCmdTrg, mergeCmdExtensions, mergeCmdName, os.Getwd)
		return err
	},
	SilenceUsage: true,
}

func init() {
	// set flag of merge command
	mergeCmd.Flags().StringVarP(&mergeCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is set as the default")
	mergeCmd.Flags().StringVarP(&mergeCmdTrg, "target", "t",
		"", "the path to the MBT results folder; the current path is set as the default")
	mergeCmd.Flags().StringSliceVarP(&mergeCmdExtensions, "extensions", "e", nil,
		"the MTA extension descriptors")
	mergeCmd.Flags().StringVarP(&mergeCmdName, "target-file-name", "n", "",
		`(required) the result file name`)
}
