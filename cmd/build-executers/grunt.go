package builders

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/cmd/constants"
	"cloud-mta-build-tool/cmd/exec"
	"cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
)

// GruntBuilder - dummy grunt builder - for testing purpose only
type GruntBuilder struct {
	path string
	name string
	dir  string
}

func (n *GruntBuilder) Path() string {
	return n.path
}

func (n *GruntBuilder) Wd() string {
	return n.dir
}

// ChangePath - to the new build path
func (n *GruntBuilder) ChangePath(newPath string) {
	n.path = newPath
}

// Build - Grunt build
func (n *GruntBuilder) Build(tdir string) error {

	logs.Logger.Debugf("Grunt builder: starting building process for path: " + tdir)
	// module path in tmp folder
	modPath := filepath.Join(tdir, n.path)
	logs.Logger.Infof("Start Building Module: " + n.name)
	// prepare Grunt commands for execution //TODO provide option to configure it from outside
	cmdParams := gruntSeq(modPath)
	// spawn build process
	err := exec.Execute(cmdParams)
	if err != nil {
		logs.Logger.Error(err)
	}
	logs.Logger.Infof("Done building module: " + n.name)
	logs.Logger.Infof("Starting archive module: " + n.name)
	// archive the module build artifacts
	err = dir.Archive(modPath, tdir+constants.DataZip, modPath)
	if err != nil {
		logs.Logger.Fatalf("Failed to archive module: " + n.name)
	}
	logs.Logger.Infof("Module %s archived successfully ", n.name)
	// After we zip the folder with the build artifacts we don't need the pre-zip folder
	// on the mtar artifacts
	err = os.RemoveAll(modPath)
	if err != nil {
		logs.Logger.Error(err)
	}
	logs.Logger.Debugf("Grunt builder: MkdirAll " + modPath)
	// Create empty folder with name as before the zip process
	err = os.MkdirAll(tdir+constants.PathSep+n.name, os.ModePerm)
	if err != nil {
		logs.Logger.Error(err)
	}

	// Move the data zip artifact to the new module folder
	err = os.Rename(tdir+constants.DataZip, tdir+constants.PathSep+n.name+constants.DataZip)
	if err != nil {
		logs.Logger.Error(err)
	}

	return nil
}

func NewGruntBuilder(p string, n string, tmpDir string) *GruntBuilder {
	return &GruntBuilder{
		path: p,
		name: n,
		dir:  tmpDir,
	}
}

func gruntSeq(modPath string) [][]string {
	cmdParams := [][]string{
		{modPath, "npm", "install"},
		{modPath, "grunt"},
		{modPath, "npm", "prune", "--production"},
	}
	return cmdParams
}
