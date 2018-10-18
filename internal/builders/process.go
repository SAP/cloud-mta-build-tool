package builders

import (
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/logs"
)

// Parse the builders command list
func Parse(data []byte) Builders {
	commands := Builders{}
	err := yaml.Unmarshal(data, &commands)
	if err != nil {
		logs.Logger.Error("Yaml file is not valid, Error: " + err.Error())
	}
	return commands
}
