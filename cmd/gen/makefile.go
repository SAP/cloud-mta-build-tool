package gen

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	"mbtv2/cmd/constants"
	fs "mbtv2/cmd/fsys"
	"mbtv2/cmd/mta/models"
	"path/filepath"
	"runtime"
	"text/template"
	"mbtv2/cmd/logs"
)

//Make - Generate the makefile
func Make() {

	var makefile = "gmake"
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

	bashFile := fs.CreateFile(projPath + constants.PathSep + makefile)
	// Read the MTA
	yamlFile, err := ioutil.ReadFile(projPath + "/mta.yaml")
	if err != nil {
		log.Printf("Not able to reay the mta.yaml file: #%v ", err)
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
		"OsCore":     OsCore,
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
	if err = t.Execute(bashFile, data); err != nil {
		logs.Logger.Error(err)
	}
	//logs.Logger.Info("MTA build script was generated successfully: " + projPath + constants.PathSep + makefile)

}

type Proc struct {
	NPROCS    string
	MAKEFLAGS string
}

// OsCore - Get the build operation's
func OsCore() []Proc {
	switch runtime.GOOS {
	case "linux":
		return []Proc{{`NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "darwin":
		return []Proc{{`NPROCS = $(shell sysctl hw.ncpu  | grep -o '[0-9]\+')`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "windows":
		return []Proc{}
	default:
		return []Proc{}
	}
}
