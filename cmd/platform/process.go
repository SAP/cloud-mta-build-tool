package platform

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func Parse(data []byte) (Platforms) {

	platforms := Platforms{}
	err := yaml.Unmarshal(data, &platforms)
	if err != nil {
		log.Printf("Yaml file is not valid, Error: " + err.Error())
		os.Exit(-1)
	}
	return platforms
}






