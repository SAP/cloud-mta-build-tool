package gen

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"text/template"
	"os"

	"mbtv2/cmd/constants"
	fs "mbtv2/cmd/fsys"
	"mbtv2/cmd/mta/models"
	"mbtv2/cmd/logs"
	"mbtv2/cmd/proc"
)

//Make - Generate the makefile
func Make() {

	var genFileName = "Makefile"
	var makeFile *os.File
	// Using the module context for the template creation
	mta := models.MTA{}
	type API map[string]string
	var data struct {
		File models.MTA
		API  API
	}
	// Get working directory
	projPath := fs.GetPath()
	// Create the init script file

	if _, err := os.Stat(projPath + constants.PathSep + genFileName); err == nil {
		// path/to/whatever exists
		makeFile = fs.CreateFile(projPath + constants.PathSep + genFileName + ".mta")
	} else {
		makeFile = fs.CreateFile(projPath + constants.PathSep + genFileName)
	}

	// Read the MTA
	yamlFile, err := ioutil.ReadFile(projPath + "/mta.yaml")
	if err != nil {
		log.Printf("Not able to read the mta.yaml file: #%v ", err)
	}
	// Parse mta
	err = yaml.Unmarshal([]byte(yamlFile), &mta)
	if err != nil {
		logs.Logger.Errorf("Not able to unmarshal the mta file ")
	}

	data.File = mta
	//// Create maps of the template method's
	funcMap := template.FuncMap{
		"ExeCommand": ExeCommand,
		"OsCore":     proc.OsCore,
	}
	// Get the path of the template source code
	_, file, _, _ := runtime.Caller(0)
	container := filepath.Join(filepath.Dir(file), "make.txt")
	// parse the template txt file
	t, err := template.New("make.txt").Funcs(funcMap).ParseFiles(container)
	if err != nil {
		panic(err)
	}
	// Execute the template
	if err = t.Execute(makeFile, data); err != nil {
		logs.Logger.Error(err)
	}
	//logs.Logger.Info("MTA build script was generated successfully: " + projPath + constants.PathSep + makefile)

}
