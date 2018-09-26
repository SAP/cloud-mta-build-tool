package tpl

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"cloud-mta-build-tool/mta"

	"cloud-mta-build-tool/cmd/builders"
	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/proc"
)

const (
	makefile       = "Makefile"
	basePre        = "base_pre.txt"
	basePost       = "base_post.txt"
	verbose        = "verbose"
	makeDefaultTpl = "make_default.txt"
	makeVerboseTpl = "make_verbose.txt"
)

func createMakeFile(path, filename string) (file *os.File, err error) {

	fullFilename := path + constants.PathSep + filename
	var mf *os.File
	if _, err = os.Stat(fullFilename); err == nil {
		mf, err = fs.CreateFile(fullFilename + ".mta")
	} else {
		mf, err = fs.CreateFile(fullFilename)
	}
	return mf, err
}

func makeFile(projectPath, makeFilename, verbTemplateName string) error {

	type API map[string]string
	// template data
	var data struct {
		File mta.MTA
		API  API
	}
	// Read the MTA
	yamlFile, err := ioutil.ReadFile(projectPath + constants.PathSep + constants.MtaYaml)
	if err != nil {
		logs.Logger.Error("Not able to read the mta file ")
		return err
	}
	mta, err := mta.Parse(yamlFile)
	if err != nil {
		logs.Logger.Error("Error occurred while parsing yaml file")
		return err
	}
	data.File = mta
	// Create maps of the template method's
	t, err := makeVerbose(verbTemplateName, basePre, basePost)
	if err != nil {
		logs.Logger.Error(err)
		return err
	}
	// Create make file for the template
	makeFile, err := createMakeFile(projectPath, makeFilename)
	if err != nil {
		logs.Logger.Error(err)
		return err
	}
	// Execute the template
	err = t.Execute(makeFile, data)
	if err != nil {
		logs.Logger.Error(err)
	}

	e := makeFile.Close()
	if err != nil {
		logs.Logger.Error(e)
	}
	return err
}

func makeVerbose(verbTemplateName string, BasePre string, BasePost string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"CommandProvider": builders.CommandProvider,
		"OsCore":          proc.OsCore,
	}
	// Get the path of the template source code
	_, file, _, _ := runtime.Caller(0)
	makeVerbPath := filepath.Join(filepath.Dir(file), verbTemplateName)
	preVerbPath := filepath.Join(filepath.Dir(file), BasePre)
	postVerbPath := filepath.Join(filepath.Dir(file), BasePost)
	// parse the template txt file
	t, err := template.New(verbTemplateName).Funcs(funcMap).ParseFiles(makeVerbPath, preVerbPath, postVerbPath)
	return t, err
}

// Make - Generate the makefile
func Make(mode string) error {
	tpl, err := makeMode(mode)
	if err != nil {
		logs.Logger.Error(err)
	}
	// Get project working directory
	pPath := fs.GetPath()
	return makeFile(pPath, makefile, tpl)
}

// Get template (default/verbose) according to the CLI flags
func makeMode(mode string) (string, error) {
	tpl := makeDefaultTpl
	if (mode == verbose) || (mode == "v") {
		tpl = makeVerboseTpl
	} else if mode == "" {
		return tpl, nil
	} else {
		return "", errors.New("command is not supported")
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
