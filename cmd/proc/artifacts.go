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
	projPath := dir.ProjectPath()
	basePath := filepath.Base(projPath)
	dir := dir.CreateDirIfNotExist(projPath + "/" + basePath)
	fmt.Print(dir)
	return dir
}
