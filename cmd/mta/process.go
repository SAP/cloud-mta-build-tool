package mta

import (
	"log"
	"cloud-mta-build-tool/cmd/mta/models"
	"gopkg.in/yaml.v2"
)

// Parse MTA file
func Parse(yamlContent []byte) (out models.MTA, err error) {
	mta := models.MTA{}
	// Format the YAML to struct's
	err = yaml.Unmarshal([]byte(yamlContent), &mta)
	if err != nil {
		log.Printf("Yaml file is not valid, Error: " + err.Error())
	}
	return mta, err
}

// Marshal - edit the file
func Marshal(in models.MTA) (mtads []byte, err error) {
	mtads, err = yaml.Marshal(&in)
	if err != nil {
		log.Printf(err.Error())
	}

	return mtads, err
}
