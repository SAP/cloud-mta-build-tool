package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"cloud-mta-build-tool/cmd/platform"
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

var pMtadSourceFlag string
var pMtadTargetFlag string

// Prepare the process for execution
var prepare = &cobra.Command{
	Use:   "prepare",
	Short: "prepare for build",
	Long:  "prepare The project generation environment For build process",
	Run: func(cmd *cobra.Command, args []string) {
		// Todo add preparation process
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
			err := packModule(args[0], args[1], args[2])
			LogError(err)
		}
	},
}

func packModule(artifactsPath string, moduleRelPath string, moduleName string) error {
	// Get module full path
	moduleFullPath, err := fs.GetFullPath(moduleRelPath)
	if err == nil {
		// Get module relative path
		moduleZipPath := filepath.Join(artifactsPath, moduleName)
		// Create empty folder with name as before the zip process
		// to put the file such as data.zip inside
		err := os.MkdirAll(moduleZipPath, os.ModePerm)
		if err != nil {
			logs.Logger.Error(err)
		} else {
			// zipping the build artifacts
			logs.Logger.Infof("Starting execute zipping module %v ", moduleName)
			moduleZipFullPath := moduleZipPath + dataZip
			if err = fs.Archive(moduleFullPath, moduleZipFullPath); err != nil {
				err = errors.New(fmt.Sprintf("Error occurred during ZIP module %v creation, error: %s  ", moduleName, err))
				removeErr := os.RemoveAll(artifactsPath)
				if removeErr != nil {
					err = errors.New(fmt.Sprintf("Error occured during directory %s removal failed %s. %s", artifactsPath, err, removeErr))
				}
			} else {
				logs.Logger.Infof("Execute zipping module %v finished successfully ", moduleName)
			}
		}
	}
	return err
}

func init() {
	pMtad.Flags().StringVarP(&pMtadSourceFlag, "source", "s", "", "Provide MTAD source ")
	pMtad.Flags().StringVarP(&pMtadTargetFlag, "target", "t", "", "Provide MTAD target ")
}

func convertTypes(mtaStr mta.MTA) {
	// Load platform configuration file
	platformCfg := platform.Parse(platform.PlatformConfig)
	// Modify MTAD object according to platform types
	// Todo platform should provided as command parameter
	platform.ConvertTypes(mtaStr, platformCfg, "cf")
}

func generateMeta(relPath string, args []string) error {
	return processMta("Metadata creation", relPath, args, func(file []byte, args []string) error {
		// Parse MTA file
		m, err := mta.ParseToMta(file)
		if err == nil {
			// Generate meta info dir with required content
			err = mta.GenMetaInfo(args[0], *m, args[1:], func(mtaStr mta.MTA) {
				convertTypes(mtaStr)
			})
		}
		return err
	})
}

// Generate metadata info from deployment
var genMeta = &cobra.Command{
	Use:   "meta",
	Short: "generate meta folder",
	Long:  "generate META-INF folder with all the required data",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := generateMeta("", args)
		LogError(err)
	},
}

func processMta(processName string, relPath string, args []string, process func(file []byte, args []string) error) error {
	logs.Logger.Info("Starting " + processName)
	s := &mta.Source{Path: relPath, Filename: "mta.yaml"}
	mf, err := s.Readfile()
	if err == nil {
		err = process(mf, args)
		if err == nil {
			logs.Logger.Info(processName + " finish successfully ")
		}
	} else {
		err = errors.New(fmt.Sprintf("MTA file not found: %s", err))
	}
	return err
}

func generateMtar(relPath string, args []string) error {
	return processMta("MTAR generation", relPath, args, func(file []byte, args []string) error {
		// Create MTAR from the building artifacts
		m, err := mta.ParseToMta(file)
		if err == nil {
			err = fs.Archive(filepath.Join(args[0]), filepath.Join(args[1], m.Id+mtarSuffix))
		}
		return err
	})
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
