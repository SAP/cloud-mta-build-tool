package gen

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"path/filepath"
	"runtime"
	"text/template"

	"cloud-mta-build-tool/cmd/builders"
	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta/models"
)

// Generate - Generate mta build file
func Generate(path string) {

	const mtaScript = "makefile"
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

	bashFile := fs.CreateFile(projPath + constants.PathSep + mtaScript)
	// Read the MTA
	yamlFile, err := ioutil.ReadFile("mta.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	// Parse mta
	err = yaml.Unmarshal([]byte(yamlFile), &mta)
	data.File = mta

	// Create maps of the template method's
	funcMap := template.FuncMap{
		"ExeCommand": builders.ExeCmd,
	}
	// Get the path of the template source code
	_, file, _, _ := runtime.Caller(0)
	container := filepath.Join(filepath.Dir(file), "script.gotmpl")
	// parse the template txt file
	t, err := template.New("script.txt").Funcs(funcMap).ParseFiles(container)
	if err != nil {
		panic(err)
	}
	// Execute the template
	if err := t.Execute(bashFile, data); err != nil {
		panic(err)
	}
	logs.Logger.Info("MTA build script was generated successfully: " + projPath + constants.PathSep + mtaScript)

}
