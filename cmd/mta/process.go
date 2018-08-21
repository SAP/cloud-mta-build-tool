package mta

import (
	"log"
	"os"

	"cloud-mta-build-tool/cmd/mta/models"
	"gopkg.in/yaml.v2"
)

// Parse MTA file
func Parse(yamlContent []byte) (out models.MTA) {
	mta := models.MTA{}
	// Format the YAML to struct's
	err := yaml.Unmarshal([]byte(yamlContent), &mta)
	if err != nil {
		log.Printf("Yaml file is not valid, Error: " + err.Error())

		os.Exit(-1)
	}
	return mta
}

// Marshal - edit the file
func Marshal(in models.MTA) []byte {
	mtads, err := yaml.Marshal(&in)
	if err != nil {
		log.Printf(err.Error())
	}

	return mtads
}
