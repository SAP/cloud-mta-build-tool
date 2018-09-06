package commands

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/cmd/mta/metainfo"
	"cloud-mta-build-tool/cmd/mta/models"
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/proc"
)

// Prepare the process for execution
var prepare = &cobra.Command{
	Use:   "prepare",
	Short: "prepare for build",
	Long:  "prepare The project generation environment For build process",
	Run: func(cmd *cobra.Command, args []string) {
		proc.Prepare()
	},
}

// Copy specific module for building purpose
var copyModule = &cobra.Command{
	Use:   "copy",
	Short: "copy module for build process",
	Long:  "copy module for build process",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Infof("Executing Copy module %v", args[2])
		proc.CopyModule(args[0], args[1])
	},
}

// Zip specific module
var pack = &cobra.Command{
	Use:   "pack",
	Short: "pack module artifacts",
	Long:  "pack the module artifacts after the build process",
	Run: func(cmd *cobra.Command, args []string) {
		// Define arguments variables
		if len(args) > 2 {
			tDir := args[0]
			mName := args[2]
			modRelPath := filepath.Join(fs.GetPath(), args[1])
			modRelName := filepath.Join(tDir, mName)
			// Create empty folder with name as before the zip process
			// to put the file such as data.zip inside
			os.MkdirAll(modRelName, os.ModePerm)
			// zipping the build artifacts
			logs.Logger.Infof("Starting execute zipping module %v ", mName)
			if err := fs.Archive(modRelPath, modRelName+constants.DataZip); err != nil {
				logs.Logger.Error("Error occurred during ZIP module %v creation, error:   ", mName, err)
				os.RemoveAll(tDir)
			} else {
				logs.Logger.Infof("Execute zipping module %v finished successfully ", mName)
			}
		} else {
			logs.Logger.Errorf("No path's provided to pack the module artifacts")
		}
	},
}

func generateMeta(relativePath string, args []string) {
	processMta(relativePath, "Metadata creation", args, func(mtaStruct models.MTA, args []string) {
		// Generate meta info dir with required content
		metainfo.GenMetaInf(args[0], mtaStruct, args[1:])
	})
}

// Generate metadata info from deployment
var genMeta = &cobra.Command{
	Use:   "meta",
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
	Run: func(cmd *cobra.Command, args []string) {
		generateMeta("", args)
	},
}

func processMta(relativePath string, processName string, args []string, process func(mtaStruct models.MTA, args []string)) {
	logs.Logger.Info("Starting " + processName)
	mtaStruct, err := proc.GetMta(filepath.Join(fs.GetPath(), relativePath))
	if err == nil {
		process(mtaStruct, args)
		logs.Logger.Info(processName + " finish successfully ")
	} else {
		logs.Logger.Error("No MTA structure found")
	}
}

func generateMtar(relativePath string, args []string) {
	processMta(relativePath, "MTAR generation", args, func(mtaStruct models.MTA, args []string) {
		// Create MTAR from the building artifacts
		fs.Archive(args[0], args[1]+constants.PathSep+mtaStruct.Id+constants.MtarSuffix)
	})
}

// Generate mtar from build artifacts
var genMtar = &cobra.Command{
	Use:   "mtar",
	Short: "generate MTAR",
	Long:  "generate MTAR from the project build artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		generateMtar("", args)
	},
}

// Cleanup temp artifacts
var cleanup = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove process artifacts",
	Long:  "Remove process artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting Cleanup process")
		// Remove temp folder
		os.RemoveAll(args[0])
		logs.Logger.Info("Done")
	},
}
