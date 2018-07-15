package proc

import (
	"fmt"
	"mbtv2/cmd/fsys"
	"path/filepath"
)

// Prepare - prepare the environment for execution
func Prepare() string {

	return mtaDir()
}

func mtaDir() string {
	projPath := dir.ProjectPath()
	basePath := filepath.Base(projPath)
	dir := dir.CreateDirIfNotExist(projPath + "/" + basePath)
	fmt.Print(dir)
	return dir
}
