// Package mta provides a convenient way of exploring the structure of `mta.yaml` file objects
// such as retrieving a list of resources required by a specific module.
package mta

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/validations"
)

// GetModules returns a list of MTA modules.
func (mta *MTA) GetModules() []*Module {
	return mta.Modules
}

// GetResources returns list of MTA resources.
func (mta *MTA) GetResources() []*Resource {
	return mta.Resources
}

// GetModuleByName returns a specific module by name.
func (mta *MTA) GetModuleByName(name string) (*Module, error) {
	for _, m := range mta.Modules {
		if m.Name == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

// GetResourceByName returns a specific resource by name.
func (mta *MTA) GetResourceByName(name string) (*Resource, error) {
	for _, r := range mta.Resources {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

// Validate validates an MTA schema.
func Validate(yamlContent []byte, projectPath string, validateSchema bool, validateProject bool) (validate.YamlValidationIssues, error) {
	//noinspection GoPreferNilSlice
	issues := []validate.YamlValidationIssue{}
	if validateSchema {
		validations, schemaValidationLog := validate.BuildValidationsFromSchemaText(schemaDef)
		if len(schemaValidationLog) > 0 {
			return schemaValidationLog, nil
		}
		yamlValidationLog, err := validate.Yaml(yamlContent, validations...)
		if err != nil && len(yamlValidationLog) == 0 {
			yamlValidationLog = append(yamlValidationLog, []validate.YamlValidationIssue{{Msg: "Validation failed" + err.Error()}}...)
		}
		issues = append(issues, yamlValidationLog...)

	}
	if validateProject {
		mtaStr := MTA{}
		Unmarshal := yaml.Unmarshal
		err := Unmarshal(yamlContent, &mtaStr)
		if err != nil {
			return nil, errors.Wrap(err, "Read failed getting MTA Yaml path reading the mta file")
		}
		projectIssues := validateYamlProject(&mtaStr, projectPath)
		issues = append(issues, projectIssues...)
	}
	return issues, nil
}

// PlatformsDefined - if platforms defined
// Only empty list of platforms indicates no platforms defined
func (module *Module) PlatformsDefined() bool {
	return module.BuildParams.SupportedPlatforms == nil || len(module.BuildParams.SupportedPlatforms) > 0
}

// Unmarshal - returns a reference to the MTA object from a byte array.
func Unmarshal(content []byte) (*MTA, error) {
	m := &MTA{}
	// Unmarshal MTA file
	err := yaml.Unmarshal([]byte(content), &m)
	if err != nil {
		err = errors.Wrap(err, "Error parsing the MTA")
	}
	return m, err
}
