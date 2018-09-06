package builders

import (
	"cloud-mta-build-tool/cmd/logs"
	"gopkg.in/yaml.v2"
)

//Parse the builders command list
func Parse(data []byte) Builders {
	commands := Builders{}
	err := yaml.Unmarshal(data, &commands)
	if err != nil {
		logs.Logger.Error("Yaml file is not valid, Error: " + err.Error())
	}
	return commands
}
