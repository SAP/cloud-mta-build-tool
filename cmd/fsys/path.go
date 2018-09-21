package dir

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/cmd/logs"
)

// GetPath - get current path
func GetPath() (dir string) {
	// TODO should get also from user
	wd, err := os.Getwd()
	if err != nil {
		logs.Logger.Panicln(err)
	}
	return wd
}

// ProjectPath - provide path for the running project
func ProjectPath() string {

	pPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Logger.Panicln(err)
	}
	return pPath
}
