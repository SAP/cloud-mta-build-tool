package tpl

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"cloud-mta-build-tool/mta"

	"cloud-mta-build-tool/cmd/builders"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/proc"
)

const (
	makefile        = "Makefile"
	basePreVerbose  = "base_pre_verbose.txt"
	basePostVerbose = "base_post_verbose.txt"
	basePreDefault  = "base_pre_default.txt"
	basePostDefault = "base_post_default.txt"
	makeDefaultTpl  = "make_default.txt"
	makeVerboseTpl  = "make_verbose.txt"
	pathSep         = string(os.PathSeparator)
)

type tplCfg struct {
	tplName string
	relPath string
	pre     string
	post    string
}

// Make - Generate the makefile
func Make(mode string) error {
	tpl, err := makeMode(mode)
	if err != nil {
		logs.Logger.Error(err)
	}
	// Get project working directory
	return makeFile(makefile, tpl)
}

func makeFile(makeFilename string, tpl tplCfg) error {

	type API map[string]string
	// template data
	var data struct {
		File mta.MTA
		API  API
	}
	// Read the MTA
	s := &mta.Source{Path: tpl.relPath}
	mf, err := s.ReadExtFile()
	// Parse
	m, e := mta.Parse(mf)
	// Template data
	data.File = m
	// Create maps of the template method's
	t, err := mapTpl(tpl.tplName, tpl.pre, tpl.post)
	if err != nil {
		logs.Logger.Error(err)
		return err
	}
	// path for creating the file
	path := filepath.Join(fs.GetPath(), tpl.relPath)
	// Create make file for the template
	makeFile, err := createMakeFile(path, makeFilename)
	if err != nil {
		logs.Logger.Error(err)
		return err
	}
	// Execute the template
	err = t.Execute(makeFile, data)
	if err != nil {
		logs.Logger.Error(err)
	}
	err = makeFile.Close()
	if err != nil {
		logs.Logger.Error(e)
	}
	return err
}

func mapTpl(templateName string, BasePre string, BasePost string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"CommandProvider": builders.CommandProvider,
		"OsCore":          proc.OsCore,
	}
	// Get the path of the template source code
	_, file, _, _ := runtime.Caller(0)
	makeVerbPath := filepath.Join(filepath.Dir(file), templateName)
	prePath := filepath.Join(filepath.Dir(file), BasePre)
	postPath := filepath.Join(filepath.Dir(file), BasePost)
	// parse the template txt file
	t, err := template.New(templateName).Funcs(funcMap).ParseFiles(makeVerbPath, prePath, postPath)
	return t, err
}

// Get template (default/verbose) according to the CLI flags
func makeMode(mode string) (tplCfg, error) {
	tpl := tplCfg{}
	if (mode == "verbose") || (mode == "v") {
		tpl.tplName = makeVerboseTpl
		tpl.pre = basePreVerbose
		tpl.post = basePostVerbose
	} else if mode == "" {
		tpl.tplName = makeDefaultTpl
		tpl.pre = basePreDefault
		tpl.post = basePostDefault
	} else {
		return tplCfg{}, errors.New("command is not supported")
	}
	return tpl, nil
}

// Find string in arg slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func createMakeFile(path, filename string) (file *os.File, err error) {

	fullFilename := path + pathSep + filename
	var mf *os.File
	if _, err = os.Stat(fullFilename); err == nil {
		mf, err = fs.CreateFile(fullFilename + ".mta")
	} else {
		mf, err = fs.CreateFile(fullFilename)
	}
	return mf, err
}
