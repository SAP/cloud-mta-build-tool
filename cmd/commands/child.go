package commands

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
)

const (
	pathSep    = string(os.PathSeparator)
	dataZip    = pathSep + "data.zip"
	mtarSuffix = ".mtar"
)

// Prepare the process for execution
var prepare = &cobra.Command{
	Use:   "prepare",
	Short: "prepare for build",
	Long:  "prepare The project generation environment For build process",
	Run: func(cmd *cobra.Command, args []string) {
		// proc.Prepare()
	},
}

// zip specific module and put the artifacts on the temp folder according
// to the mtar structure, i.e each module have new entry as folder in the mtar folder
// Note - even if the path of the module was changed in the mta.yaml in the mtar the
// the module folder will get the module name
var pack = &cobra.Command{
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
			packModule(args[0], args[1], args[2])
		}
	},
}

func packModule(tDir string, mPathProp string, mNameProp string) {
	// Get module path
	mp := filepath.Join(fs.GetPath(), mPathProp)
	// Get module relative path
	mrp := filepath.Join(tDir, mNameProp)
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(mrp, os.ModePerm)
	if err != nil {
		logs.Logger.Error(err)
	} else {
		// zipping the build artifacts
		logs.Logger.Infof("Starting execute zipping module %v ", mNameProp)
		if err = fs.Archive(mp, mrp+dataZip); err != nil {
			logs.Logger.Errorf("Error occurred during ZIP module %v creation, error: %s  ", mNameProp, err)
			err = os.RemoveAll(tDir)
			if err != nil {
				logs.Logger.Error(err)
			}
		} else {
			logs.Logger.Infof("Execute zipping module %v finished successfully ", mNameProp)
		}
	}
}

func generateMeta(relPath string, args []string) {
	processMta("Metadata creation", relPath, args, func(file []byte, args []string) {
		m, err := mta.Parse(file)
		if err != nil {
			logs.Logger.Error(err)
		}
		// Generate meta info dir with required content
		mta.GenMetaInfo(args[0], m, args[1:])
	})
}

// Generate metadata info from deployment
var genMeta = &cobra.Command{
	Use:   "meta",
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		generateMeta("", args)
	},
}

func processMta(processName string, relPath string, args []string, process func(file []byte, args []string)) {
	logs.Logger.Info("Starting " + processName)
	mf, err := getFile(relPath)
	if err == nil {
		process(mf, args)
		logs.Logger.Info(processName + " finish successfully ")
	} else {
		logs.Logger.Error("MTA file not found")
	}
}

func generateMtar(relPath string, args []string) error {
	processMta("MTAR generation", relPath, args, func(file []byte, args []string) {
		// Create MTAR from the building artifacts
		m, err := mta.Parse(file)
		if err != nil {
			logs.Logger.Error(err)
		}
		err = fs.Archive(args[0], args[1]+pathSep+m.Id+mtarSuffix)
		if err != nil {
			logs.Logger.Error(err)
		}
	})
	return nil
}

// Generate mtar from build artifacts
var genMtar = &cobra.Command{
	Use:   "mtar",
	Short: "generate MTAR",
	Long:  "generate MTAR from the project build artifacts",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		generateMtar("", args)
	},
}

// Cleanup temp artifacts
var cleanup = &cobra.Command{
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
