package proc

import (
	"cloud-mta-build-tool/cmd/fsys"
	"fmt"
	"path/filepath"
)

// Prepare - prepare the environment for execution
func Prepare() string {

	return mtaDir()
}

//Todo should be part of the MakeFile
func mtaDir() string {
	projPath := dir.ProjectPath()
	basePath := filepath.Base(projPath)
	dir := dir.CreateDirIfNotExist(projPath + "/" + basePath)
	fmt.Print(dir)
	return dir
}
