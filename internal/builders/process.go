package builders

import (
	"gopkg.in/yaml.v2"
)

// parse the builders command list
func parse(data []byte) (Builders, error) {
	commands := Builders{}
	err := yaml.Unmarshal(data, &commands)
	if err != nil {
		return Builders{}, err
	}
	return commands, nil
}
