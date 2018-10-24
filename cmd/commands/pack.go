package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	fs "cloud-mta-build-tool/internal/fsys"

	"cloud-mta-build-tool/internal/logs"
)

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
