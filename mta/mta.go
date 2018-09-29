package mta

import (
	"fmt"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/logs"
)

// MTA - Main mta struct
type MTA struct {
	SchemaVersion *string      `yaml:"_schema-version"`
	Id            string       `yaml:"ID"`
	Version       string       `yaml:"version,omitempty"`
	Modules       []*Modules   `yaml:"modules,omitempty"`
	Resources     []*Resources `yaml:"resources,omitempty"`
	Parameters    Parameters   `yaml:"parameters,omitempty"`
}

// Parse MTA file and provide mta object with data
func Parse(yamlContent []byte) (out MTA, err error) {
	mta := MTA{}
	// Format the YAML to struct's
	err = yaml.Unmarshal([]byte(yamlContent), &mta)
	if err != nil {
		logs.Logger.Error("Yaml file is not valid, Error: " + err.Error())
	}
	return mta, err
}

// Marshal - usage for edit purpose
func Marshal(in MTA) (mtads []byte, err error) {
	mtads, err = yaml.Marshal(&in)
	if err != nil {
		logs.Logger.Error(err.Error())
	}
	return mtads, err
}

// Provide list of modules
func modules(mta MTA) []string {
	var mNames []string
	for _, mod := range mta.Modules {
		mNames = append(mNames, mod.Name)
	}
	return mNames
}

// GetModulesNames - get list of modules names
func GetModulesNames(file []byte) ([]string, error) {

	mta, err := Parse(file)
	if err != nil {
		return nil, fmt.Errorf("not able to read the mta file : %s", err.Error())
	}
	return modules(mta), err
}
