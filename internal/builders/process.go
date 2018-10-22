package builders

import (
	"gopkg.in/yaml.v2"
)

// Parse the builders command list
func Parse(data []byte) (Builders , error){
	commands := Builders{}
	err := yaml.Unmarshal(data, &commands)
	if err != nil {
		return Builders{}, err
	}
	return commands,nil
}
