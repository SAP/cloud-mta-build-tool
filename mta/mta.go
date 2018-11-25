// Package mta provides a convenient way of exploring the structure of `mta.yaml` file objects
// such as retrieving a list of resources required by a specific module.
package mta

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/validations"
)

// Unmarshal un-marshals an MTA byte document and provides an MTA object with corresponding.
func (mta *MTA) Unmarshal(yamlContent []byte) (err error) {
	// Format the YAML to struct's
	err = yaml.Unmarshal([]byte(yamlContent), &mta)
	if err != nil {
		return errors.Wrap(err, "error occurred while parsing file : %s")
	}
	return nil
}

// Marshal serializes the MTA into an encoded YAML document.
func Marshal(in *MTA) (mtads []byte, err error) {
	mtads, err = yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	return mtads, nil
}

// Read returns mta byte slice.
func Read(ep *Loc) ([]byte, error) {
	fileFullPath, err := ep.GetMtaYamlPath()
	if err != nil {
		return nil, errors.Wrap(err, "Read failed getting MTA Yaml path")
	}
	// ParseFile MTA file
	yamlFile, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading the MTA file")
	}
	return yamlFile, nil
}

// ParseFile returns a reference to the MTA object from a given mta.yaml file.
func ParseFile(ep *Loc) (*MTA, error) {
	var mta *MTA
	yamlContent, err := Read(ep)
	// ParseFile MTA file
	if err == nil {
		mta, err = ParseByte(yamlContent)
	}
	return mta, err
}

// ParseByte returns a reference to the MTA object from a byte array.
func ParseByte(content []byte) (*MTA, error) {
	mta := &MTA{}
	// Unmarshal MTA file
	err := mta.Unmarshal(content)
	if err != nil {
		err = errors.Wrap(err, "Error parsing the MTA")
	}
	return mta, err
}

// GetModules returns a list of MTA modules.
func (mta *MTA) GetModules() []*Modules {
	return mta.Modules
}

// GetResources returns list of MTA resources.
func (mta *MTA) GetResources() []*Resources {
	return mta.Resources
}

// GetModuleByName returns a specific module by name.
func (mta *MTA) GetModuleByName(name string) (*Modules, error) {
	for _, m := range mta.Modules {
		if m.Name == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

// GetResourceByName returns a specific resource by name.
func (mta *MTA) GetResourceByName(name string) (*Resources, error) {
	for _, r := range mta.Resources {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

// GetModulesNames returns a list of module names.
func (mta *MTA) GetModulesNames() ([]string, error) {
	return mta.getModulesOrder()
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
		mta := MTA{}
		Unmarshal := yaml.Unmarshal
		err := Unmarshal(yamlContent, &mta)
		if err != nil {
			return nil, errors.Wrap(err, "Read failed getting MTA Yaml path reading the mta file")
		}
		projectIssues := validateYamlProject(&mta, projectPath)
		issues = append(issues, projectIssues...)
	}
	return issues, nil
}

// PlatformsDefined - if platforms defined
// Only empty list of platforms indicates no platforms defined
func (module *Modules) PlatformsDefined() bool {
	return module.BuildParams.SupportedPlatforms == nil || len(module.BuildParams.SupportedPlatforms) > 0
}
