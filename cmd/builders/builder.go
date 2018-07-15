package builders

import (
	"mbtv2/cmd/constants"
	"mbtv2/cmd/fsys"
	"mbtv2/cmd/logs"
)

// Builder - A builder is what used to build the language. It should be able to change working dir.
type Builder interface {
	Build(path string) error   // Build builds the code at current dir. It returns an error if failed.
	Path() string              // Path returns the current working dir.
	ChangePath(newPath string) // ChangePath changes the working dir to newPath.
	Wd() string                // Current working dir
}

// TempDirFunc is what generates a new temp dir. Golang would requires it in GOPATH, so make it changeable.
type TempDirFunc func() string

// Build - Generic Build
func Build(b Builder, toPath string, mkTempDir string) error {

	logs.Logger.Debugf("Base builder:path: " + toPath)
	// TODO support build temp target for each module
	//if mkTempDir == nil {
	//	mkTempDir()
	//}

	// Get module path
	path := b.Path()
	tmpWd := b.Wd()
	// TODO -Remove the directory after the function exit
	// Copy module from source to target project -> tmp dir
	if err := dir.CopyDir(toPath+constants.PathSep+path, tmpWd+constants.PathSep+path); err != nil {
		logs.Logger.Fatalf("Base builder: CopyDir:err: ", err)
		return err
	}
	logs.Logger.Debugf("Base builder: tmpDir path: " + tmpWd)
	// Build specific module artifacts
	err := b.Build(tmpWd)
	if err != nil {
		logs.Logger.Fatalf("Base builder: Build:err: ", err)
		return err
	}

	return nil
}
