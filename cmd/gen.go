package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

// mtad command flags
var mtadCmdSrc string
var mtadCmdTrg string
var mtadCmdExtensions []string
var mtadCmdPlatform string

// meta command flags
var metaCmdSrc string
var metaCmdTrg string
var metaCmdDesc string
var metaCmdExtensions []string
var metaCmdPlatform string

// mtar command flags
var mtarCmdSrc string
var mtarCmdTrg string
var mtarCmdDesc string
var mtarCmdTrgProvided string
var mtarCmdExtensions []string
var mtarCmdMtarName string

// init - inits flags of init command
func init() {

	// set flags of mtad command
	mtadCmd.Flags().StringVarP(&mtadCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is set as the default")
	mtadCmd.Flags().StringVarP(&mtadCmdTrg, "target", "t",
		"", "the path to the MBT results folder; the current path is set as the default")
	mtadCmd.Flags().StringSliceVarP(&mtadCmdExtensions, "extensions", "e", nil,
		"the MTA extension descriptors")
	mtadCmd.Flags().StringVarP(&mtadCmdPlatform, "platform", "p", "cf",
		`the deployment platform; supported platforms: "cf" (default value), "xsa", "neo"`)
	mtadCmd.Flags().BoolP("help", "h", false, `prints detailed information about the "mtad" command`)

	// set flags of meta command
	metaCmd.Flags().StringVarP(&metaCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is set as the default")
	metaCmd.Flags().StringVarP(&metaCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is set as the default")
	metaCmd.Flags().StringVarP(&metaCmdDesc, "desc", "d", "",
		`the MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	metaCmd.Flags().StringSliceVarP(&metaCmdExtensions, "extensions", "e", nil,
		"the MTA extension descriptors")
	metaCmd.Flags().StringVarP(&metaCmdPlatform, "platform", "p", "cf",
		`the deployment platform; supported platforms: "cf" (default value), "xsa", "neo"`)
	metaCmd.Flags().BoolP("help", "h", false, `prints detailed information about the "meta" command`)

	// set flags of mtar command
	mtarCmd.Flags().StringVarP(&mtarCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is set as the default")
	mtarCmd.Flags().StringVarP(&mtarCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is set as the default")
	mtarCmd.Flags().StringVarP(&mtarCmdDesc, "desc", "d", "",
		`the MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	mtarCmd.Flags().StringSliceVarP(&mtarCmdExtensions, "extensions", "e", nil,
		"the MTA extension descriptors")
	mtarCmd.Flags().StringVarP(&mtarCmdMtarName, "mtar", "m", "*",
		"the archive name")
	mtarCmd.Flags().StringVarP(&mtarCmdTrgProvided, "target_provided", "", "",
		"the MTA target provided indicator; supported values: true, false")
	_ = mtarCmd.Flags().MarkHidden("target_provided")
	mtarCmd.Flags().BoolP("help", "h", false, `prints detailed information about the "mtar" command`)

}

// Provide mtad.yaml from mta.yaml
var mtadCmd = &cobra.Command{
	Use:   "mtad",
	Short: "Generates MTAD",
	Long:  "Generates deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteGenMtad(mtadCmdSrc, mtadCmdTrg, mtadCmdExtensions, mtadCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Generate metadata info from deployment
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Generates the META-INF folder",
	Long:  "Generates META-INF folder with manifest and MTAD files",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteGenMeta(metaCmdSrc, metaCmdTrg, metaCmdDesc, metaCmdExtensions, metaCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Generate mtar from build artifacts
var mtarCmd = &cobra.Command{
	Use:   "mtar",
	Short: "Generates MTA archive",
	Long:  "Generates MTA archive from the folder with all artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteGenMtar(mtarCmdSrc, mtarCmdTrg, mtarCmdTrgProvided, mtarCmdDesc, mtadCmdExtensions, mtarCmdMtarName, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}
