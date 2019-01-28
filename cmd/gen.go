package commands

import (
	"os"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/artifacts"
)

// mtad command flags
var mtadCmdSrc string
var mtadCmdTrg string
var mtadCmdDesc string
var mtadCmdPlatform string

// meta command flags
var metaCmdSrc string
var metaCmdTrg string
var metaCmdDesc string
var metaCmdPlatform string

// mtar command flags
var mtarCmdSrc string
var mtarCmdTrg string
var mtarCmdDesc string

// init - inits flags of init command
func init() {

	// set flags of mtad command
	mtadCmd.Flags().StringVarP(&mtadCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is default")
	mtadCmd.Flags().StringVarP(&mtadCmdTrg, "target", "t",
		"", "the path to the MBT results folder; the current path is default")
	mtadCmd.Flags().StringVarP(&mtadCmdDesc, "desc", "d", "",
		"the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")
	mtadCmd.Flags().StringVarP(&mtadCmdPlatform, "platform", "p", "", "Provide MTA platform ")

	// set flags of meta command
	metaCmd.Flags().StringVarP(&metaCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is default")
	metaCmd.Flags().StringVarP(&metaCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is default")
	metaCmd.Flags().StringVarP(&metaCmdDesc, "desc", "d", "",
		"the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")
	metaCmd.Flags().StringVarP(&metaCmdPlatform, "platform", "p", "",
		"the deployment platform; supported plaforms: cf, xsa")

	// set flags of mtar command
	mtarCmd.Flags().StringVarP(&mtarCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is default")
	mtarCmd.Flags().StringVarP(&mtarCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is default")
	mtarCmd.Flags().StringVarP(&mtarCmdDesc, "desc", "d", "",
		"the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")
}

// Provide mtad.yaml from mta.yaml
var mtadCmd = &cobra.Command{
	Use:   "mtad",
	Short: "generates MTAD",
	Long:  "generates deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteGenMtad(mtadCmdSrc, mtadCmdTrg, mtadCmdDesc, mtadCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}

// Generate metadata info from deployment
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "generates the META-INF folder",
	Long:  "generates META-INF folder with manifest and MTAD files",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteGenMeta(metaCmdSrc, metaCmdTrg, metaCmdDesc, metaCmdPlatform, true, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}

// Generate mtar from build artifacts
var mtarCmd = &cobra.Command{
	Use:   "mtar",
	Short: "generates MTA archive",
	Long:  "generates MTA archive from the folder with all artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteGenMtar(mtarCmdSrc, mtarCmdTrg, mtarCmdDesc, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}
