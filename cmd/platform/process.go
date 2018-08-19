package platform

import (
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/logs"
)

func Parse(data []byte) (Platforms) {

	platforms := Platforms{}
	err := yaml.Unmarshal(data, &platforms)
	if err != nil {
		logs.Logger.Error("Yaml file is not valid, Error: " + err.Error())
	}
	return platforms
}
