package commands

import (
	"gopkg.in/yaml.v2"
)

// parse the builders command list
func parseBuilders(data []byte) (Builders, error) {
	builders := Builders{}
	err := yaml.Unmarshal(data, &builders)
	if err != nil {
		return Builders{}, err
	}
	return builders, nil
}

// parse the module types
func parseModuleTypes(data []byte) (ModuleTypes, error) {
	moduleTypes := ModuleTypes{}
	err := yaml.Unmarshal(data, &moduleTypes)
	if err != nil {
		return ModuleTypes{}, err
	}
	return moduleTypes, nil
}
