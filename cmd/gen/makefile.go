package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/builders"
	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta/models"
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
	mta := models.MTA{}
	type API map[string]string
	var data struct {
		File models.MTA
		API  API
	}
	// Read the MTA
	yamlFile, err := ioutil.ReadFile(yamlPath + "/" + yamlName)

	if err != nil {
		logs.Logger.Error("Not able to read the mta file ")
		return err
	}
	// Parse mta
	err = yaml.Unmarshal([]byte(yamlFile), &mta)
	if err != nil {
		logs.Logger.Error("Not able to unmarshal the mta file ")
		return err
	}

	data.File = mta
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
	//logs.Logger.Info("MTA build script was generated successfully: " + projPath + constants.PathSep + makefile)
}

//Make - Generate the makefile
func Make() error {

	const MakeVerbTmpl = "make_verbose.txt"
	var genFileName = "Makefile"

	// Get working directory
	projPath := fs.GetPath()

	return makeFile(projPath, "mta.yaml", projPath, genFileName, MakeVerbTmpl)

}
