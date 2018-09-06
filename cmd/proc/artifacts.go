package proc

import (
	"fmt"
	"path/filepath"

	"cloud-mta-build-tool/cmd/fsys"
)

// Prepare - future use pre-process - prepare the environment for execution
func Prepare() string {
	return mtaDir()
}

func mtaDir() string {
	pPath := dir.ProjectPath()
	basePath := filepath.Base(pPath)
	dirName := pPath + "/" + basePath
	dir.CreateDirIfNotExist(dirName)
	fmt.Print(dirName)
	return dirName
}
