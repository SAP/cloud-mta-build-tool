package builders

import (
	"log"
	"os"
	"path/filepath"

	"mbtv2/cmd/constants"
	"mbtv2/cmd/exec"
	"mbtv2/cmd/fsys"
	"mbtv2/cmd/logs"
)

type NpmBuilder struct {
	path string
	name string
	dir  string
}

// Path - Provide the path
func (n *NpmBuilder) Path() string {
	return n.path
}

// ChangePath - Change the  module path
func (n *NpmBuilder) ChangePath(newPath string) {
	logs.Logger.Debugf("Build: NodeBuilder: changing to ", newPath)
	n.path = newPath

}

// Wd - working dir
func (n *NpmBuilder) Wd() string {
	return n.dir
}

// TempDir - get temp dir
func (n *NpmBuilder) TempDir() string {
	return n.dir
}

// Build - Build node module
func (n *NpmBuilder) Build(pdir string) error {

	log.Println("Start Building Module " + n.name)
	// module Path
	modPath := filepath.Join(pdir, n.path)
	logs.Logger.Println("Start Building Module: " + n.name)
	// prepare npm commands for execution
	cmdParams := npmSeq(modPath)
	// spawn build process
	exec.Execute(cmdParams)
	logs.Logger.Println("Done building module: " + n.name)
	logs.Logger.Info("Starting archive modules: " + n.name)
	// archive the module build artifacts
	err := dir.Archive(modPath, pdir+constants.PathSep+constants.DataZip, modPath)
	if err != nil {
		logs.Logger.Error("Failed to archive module: " + n.name)
	}
	logs.Logger.Infof("Module %s archived successfully ", n.name)
	// Remove the zipped dir
	os.RemoveAll(modPath)
	os.MkdirAll(modPath, os.ModePerm)
	os.Rename(pdir+constants.PathSep+constants.DataZip, modPath+constants.PathSep+constants.DataZip)

	return nil

}

// NewNPMBuilder - npm builder instance
func NewNPMBuilder(p string, n string, tmpDir string) *NpmBuilder {
	return &NpmBuilder{
		path: p,
		name: n,
		dir:  tmpDir,
	}
}

func npmSeq(modPath string) [][]string {
	cmdParams := [][]string{
		{modPath, "npm", "install", "--production"},
	}
	return cmdParams
}
