package commands

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta/metainfo"
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

// zip specific module and put the artifacts on the temp folder according
// to the mtar structure, i.e each module have new entry as folder in the mtar folder
// note - even if the path of the module was changed in the mta.yaml in the mtar the
// the module folder will get the module name
var pack = &cobra.Command{
	Use:   "pack",
	Short: "pack module artifacts",
	Long:  "pack the module artifacts after the build process",
	Run: func(cmd *cobra.Command, args []string) {
		// Define arguments variables
		if len(args) > 0 {
			// path of the temp directory to add the build module
			td := args[0]
			mName := args[2]
			mRelPath := fs.ProjectPath() + "/" + args[1]
			modRelName := filepath.Join(td, mName)
			// Create empty folder with name as before the zip process
			// to put the file such as data.zip inside
			err := os.MkdirAll(modRelName, os.ModePerm)
			if err != nil {
				logs.Logger.Error(err)
			}
			// zipping the build artifacts
			logs.Logger.Infof("Starting execute zipping module %v ", mName)
			if err := fs.Archive(mRelPath, td+"/"+args[2]+constants.DataZip, mRelPath); err != nil {
				logs.Logger.Error("Error occurred during ZIP module %v creation, error:   ", args[0], err)
			}
			logs.Logger.Infof("Execute zipping module %v finished successfully ", mName)
		} else {
			logs.Logger.Errorf("No path's provided to pack the module artifacts")
		}
	},
}

// Generate metadata info from deployment
var genMeta = &cobra.Command{
	Use:   "meta",
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting execute metadata creation")
		ms, err := proc.GetMta(fs.GetPath())
		if err != nil {
			logs.Logger.Error(err)
		}
		mtarDir := args[0]
		// Generate meta info dir with required content
		metainfo.GenMetaInf(mtarDir, ms, args[1:])
		logs.Logger.Info("Metadata creation finish successfully ")
	},
}

// Generate mtar from build artifacts
var genMtar = &cobra.Command{
	Use:   "mtar",
	Short: "generate MTAR",
	Long:  "generate MTAR from the project build artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting execute Build of MTAR")
		ms, err := proc.GetMta(fs.GetPath())
		if err != nil {
			logs.Logger.Error(err)
		}
		tDir := args[0]
		pDir := args[1]
		// Create MTAR from the building artifacts
		err = fs.Archive(tDir, pDir+constants.PathSep+ms.Id+constants.MtarSuffix, tDir)
		if err != nil {
			logs.Logger.Error(err)
		}
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
		err := os.RemoveAll(args[0])
		if err != nil {
			logs.Logger.Error(err)
		}
		logs.Logger.Info("Done")
	},
}
