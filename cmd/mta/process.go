package mta

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"io/ioutil"

	"mbtv2/cmd/logs"
	"mbtv2/cmd/mta/models"
)

// Parse MTA file
func Parse(yamlContent []byte) (out models.MTA, err error) {
	mta := models.MTA{}
	// Format the YAML to struct's
	err = yaml.Unmarshal([]byte(yamlContent), &mta)
	if err != nil {
		log.Printf("Yaml file is not valid, Error: " + err.Error())

		os.Exit(-1)
	}
	return mta, err
}

// Marshal - edit the file
func Marshal(in models.MTA) []byte {
	mtads, err := yaml.Marshal(&in)
	if err != nil {
		log.Printf(err.Error())
	}

	return mtads
}

// Load - load the mta.yaml file
func Load(path string) []byte {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Logger.Errorf("mta.yaml not found for path %s, error is: #%v ", path, err)
		// YAML descriptor file not found abort the process
		os.Exit(1)
	}
	logs.Logger.Debugf("The file loaded successfully:" + string(yamlFile))
	return yamlFile
}
