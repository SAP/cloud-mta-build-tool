package builders

import (
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/logs"
)

//Parse the builders command list
func Parse(data []byte) (ExeCommands) {
	commands := ExeCommands{}
	err := yaml.Unmarshal(data, &commands)
	if err != nil {
		logs.Logger.Error("Yaml file is not valid, Error: " + err.Error())
	}
	return commands
}
