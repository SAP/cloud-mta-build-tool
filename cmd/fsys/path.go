package dir

import (
	"cloud-mta-build-tool/cmd/logs"
	"os"
	"path/filepath"
)

func GetPath() (dir string) {
	// TODO should get also from user
	wd, err := os.Getwd()
	if err != nil {
		logs.Logger.Panicln(err)
	}
	return wd
}

func ProjectPath() string {

	projPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Logger.Panicln(err)
	}
	return projPath
}
