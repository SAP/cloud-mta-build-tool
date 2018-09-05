package proc

import (
	"fmt"
	"path/filepath"

	"cloud-mta-build-tool/cmd/fsys"
)

// Prepare - prepare the environment for execution
func Prepare() string {
	return mtaDir()
}

//Todo should be part of the MakeFile (mkdir) , will support generic pre build process
func mtaDir() string {
	pPath := dir.ProjectPath()
	basePath := filepath.Base(pPath)
	dirName := pPath + "/" + basePath
	dir.CreateDirIfNotExist(dirName)
	fmt.Print(dirName)
	return dirName
}
