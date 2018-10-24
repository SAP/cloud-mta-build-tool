package commands

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
)

const (
	pathSep    = string(os.PathSeparator)
	dataZip    = pathSep + "data.zip"
	mtarSuffix = ".mtar"
)

var pMtadSourceFlag string
var pMtadTargetFlag string

func init() {
	genMtadCmd.Flags().StringVarP(&pMtadSourceFlag, "source", "s", "", "Provide MTAD source ")
	genMtadCmd.Flags().StringVarP(&pMtadTargetFlag, "target", "t", "", "Provide MTAD target ")
}

// zip specific module and put the artifacts on the temp folder according
// to the mtar structure, i.e each module have new entry as folder in the mtar folder
// Note - even if the path of the module was changed in the mta.yaml in the mtar the
// the module folder will get the module name
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack module artifacts",
	Long:  "pack the module artifacts after the build process",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return nil
		} else {
			return errors.New("no path's provided to pack the module artifacts")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			err := packModule(args[0], args[1], args[2])
			LogError(err)
		}
	},
}

// Generate metadata info from deployment
var genMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := generateMeta("", args)
		LogError(err)
	},
}

// Generate mtar from build artifacts
var genMtarCmd = &cobra.Command{
	Use:   "mtar",
	Short: "generate MTAR",
	Long:  "generate MTAR from the project build artifacts",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		generateMtar("", args)
	},
}

// Provide mtad.yaml from mta.yaml
var genMtadCmd = &cobra.Command{
	Use:   "mtad",
	Short: "Provide mtad",
	Long:  "Provide deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		mtaStr, err := mta.ReadMta(pMtadSourceFlag, "mta.yaml")
		if err == nil {
			err = mta.GenMtad(*mtaStr, pMtadTargetFlag, func(mtaStr mta.MTA) {
				convertTypes(mtaStr)
			})
		}
		LogError(err)
	},
}

// Validate mta.yaml
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "MBT validation",
	Long:  "MBT validation process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		validateSchema, validateProject, err := getValidationMode(validationFlag)
		if err == nil {
			err = validateMtaYaml("", "mta.yaml", validateSchema, validateProject)
		}
		LogError(err)
	},
}

// Cleanup temp artifacts
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove process artifacts",
	Long:  "Remove MTA build process artifacts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting Cleanup process")
		// Remove temp folder
		err := os.RemoveAll(args[0])
		if err != nil {
			logs.Logger.Error(err)
		} else {
			logs.Logger.Info("Done")
		}
	},
}
