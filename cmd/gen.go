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
	mtadCmd.Flags().StringVarP(&mtadCmdSrc, "source", "s", "", "Provide MTA source ")
	mtadCmd.Flags().StringVarP(&mtadCmdTrg, "target", "t", "", "Provide MTA target ")
	mtadCmd.Flags().StringVarP(&mtadCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	mtadCmd.Flags().StringVarP(&mtadCmdPlatform, "platform", "p", "", "Provide MTA platform ")

	// set flags of meta command
	metaCmd.Flags().StringVarP(&metaCmdSrc, "source", "s", "", "Provide MTA source ")
	metaCmd.Flags().StringVarP(&metaCmdTrg, "target", "t", "", "Provide MTA target ")
	metaCmd.Flags().StringVarP(&metaCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	metaCmd.Flags().StringVarP(&metaCmdPlatform, "platform", "p", "", "Provide MTA platform ")

	// set flags of mtar command
	mtarCmd.Flags().StringVarP(&mtarCmdSrc, "source", "s", "", "Provide MTA source ")
	mtarCmd.Flags().StringVarP(&mtarCmdTrg, "target", "t", "", "Provide MTA target ")
	mtarCmd.Flags().StringVarP(&mtarCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
}

// Provide mtad.yaml from mta.yaml
var mtadCmd = &cobra.Command{
	Use:   "mtad",
	Short: "Provide mtad",
	Long:  "Provide deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
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
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
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
	Short: "generate MTAR",
	Long:  "generate MTAR from the project build artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteGenMtar(mtarCmdSrc, mtarCmdTrg, mtarCmdDesc, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}
