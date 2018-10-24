package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"

	"cloud-mta-build-tool/internal/platform"
	"cloud-mta-build-tool/mta"
)

// generate build metadata artifacts
func generateMeta(relPath string, args []string) error {
	return processMta("Metadata creation", relPath, args, func(file []byte, args []string) error {
		// Parse MTA file
		m, err := mta.ParseToMta(file)
		if err == nil {
			// Generate meta info dir with required content
			err = mta.GenMetaInfo(args[0], *m, args[1:], func(mtaStr mta.MTA) {
				err = convertTypes(mtaStr)
			})
		}
		return err
	})
}

// generate mtar archive from the build artifacts
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

// convert types to appropriate target platform types
func convertTypes(mtaStr mta.MTA) error {
	// Load platform configuration file
	platformCfg, err := platform.Parse(platform.PlatformConfig)
	if err == nil {
		// Modify MTAD object according to platform types
		// Todo platform should provided as command parameter
		platform.ConvertTypes(mtaStr, platformCfg, "cf")
	}
	return err
}

// process mta.yaml file
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

// pack build module artifacts
func packModule(artifactsPath string, moduleRelPath string, moduleName string) error {
	// Get module full path
	moduleFullPath, err := fs.GetFullPath(moduleRelPath)
	if err == nil {
		// Get module relative path
		moduleZipPath := filepath.Join(artifactsPath, moduleName)
		// Create empty folder with name as before the zip process
		// to put the file such as data.zip inside
		err = os.MkdirAll(moduleZipPath, os.ModePerm)
		if err == nil {
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
