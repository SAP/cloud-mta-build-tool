package tpl

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"cloud-mta-build-tool/mta"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/builders"
	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/proc"
)

func createMakeFile(path, filename string) (file *os.File, err error) {

	fullFilename := path + constants.PathSep + filename

	var mf *os.File
	if _, sErr := os.Stat(fullFilename); sErr == nil {
		// path/to/whatever exists
		mf, sErr = fs.CreateFile(fullFilename + ".mta")
	} else {
		mf, sErr = fs.CreateFile(fullFilename)
	}
	return mf, err
}

func makeFile(yamlPath, yamlName, makeFilePath, makeFilename, verbTemplateName string) error {

	const BasePre = "base_pre.txt"
	const BasePost = "base_post.txt"

	// Using the module context for the template creation
	m := mta.MTA{}
	type API map[string]string
	var data struct {
		File mta.MTA
		API  API
	}
	// Read the MTA
	yamlFile, err := ioutil.ReadFile(yamlPath + constants.PathSep + yamlName)

	if err != nil {
		logs.Logger.Error("Not able to read the mta file ")
		return err
	}
	// Parse mta
	err = yaml.Unmarshal([]byte(yamlFile), &m)
	if err != nil {
		logs.Logger.Error("Not able to unmarshal the mta file ")
		return err
	}

	data.File = m
	// Create maps of the template method's
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
	if err != nil {
		logs.Logger.Error(err)
		return err
	}
	makeFile, err := createMakeFile(makeFilePath, makeFilename)
	if err != nil {
		logs.Logger.Error(err)
		return err
	}

	// Execute the template
	err = t.Execute(makeFile, data)
	if err != nil {
		logs.Logger.Error(err)
	}

	makeFile.Close()
	return err
}

// Make - Generate the makefile
func Make(mode []string) error {
	tpl, err := makeMode(mode)
	if err != nil {
		logs.Logger.Error(err)
	}
	var genFileName = "Makefile"
	// Get working directory
	projPath := fs.GetPath()
	return makeFile(projPath, "mta.yaml", projPath, genFileName, tpl)

}

// Get template according to the CLI flags
func makeMode(mode []string) (string, error) {
	tpl := "make_default.txt"
	if (len(mode) > 0) && (stringInSlice("--verbose", mode)) {
		if (mode[0] == "--verbose") || (mode[0] == "-v") {
			tpl = "make_verbose.txt"
		}
	} else if len(mode) == 0 {
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
