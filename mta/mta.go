package mta

import (
	"fmt"

	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/validations"
)

type MTAI interface {
	GetModules() []*Modules
	GetResources() []*Resources
	GetModuleByName(name string) (*Modules, error)
	GetModulesNames() []string
	GetResourceByName(name string) (*Resources, error)
}

// MTA - Main mta struct
type MTA struct {
	SchemaVersion *string      `yaml:"_schema-version"`
	Id            string       `yaml:"ID"`
	Version       string       `yaml:"version,omitempty"`
	Modules       []*Modules   `yaml:"modules,omitempty"`
	Resources     []*Resources `yaml:"resources,omitempty"`
	Parameters    Parameters   `yaml:"parameters,omitempty"`
}

// BuildParameters - build params
type BuildParameters struct {
	Builder  string          `yaml:"builder,omitempty"`
	Type     string          `yaml:"type,omitempty"`
	Path     string          `yaml:"path,omitempty"`
	Requires []BuildRequires `yaml:"requires,omitempty"`
}

// Modules - MTA modules
type Modules struct {
	Name        string
	Type        string
	Path        string          `yaml:"path,omitempty"`
	Requires    []Requires      `yaml:"requires,omitempty"`
	Provides    []Provides      `yaml:"provides,omitempty"`
	Parameters  Parameters      `yaml:"parameters,omitempty"`
	BuildParams BuildParameters `yaml:"build-parameters,omitempty"`
	Properties  Properties      `yaml:"properties,omitempty"`
}

// Properties - MTA map
type Properties map[string]interface{}

// Parameters - MTA parameters
type Parameters map[string]interface{}

// Provides - MTA struct
type Provides struct {
	Name       string
	Properties Properties `yaml:"properties,omitempty"`
}

// Requires / Mta struct
type Requires struct {
	Name       string     `yaml:"name,omitempty"`
	Group      string     `yaml:"group,omitempty"`
	Type       string     `yaml:"type,omitempty"`
	Properties Properties `yaml:"properties,omitempty"`
}

// BuildRequires - build requires section
type BuildRequires struct {
	Name       string `yaml:"name,omitempty"`
	TargetPath string `yaml:"target-path,omitempty"`
}

// Resources - resources section
type Resources struct {
	Name       string
	Type       string
	Parameters Parameters `yaml:"parameters,omitempty"`
	Properties Properties `yaml:"properties,omitempty"`
}

type MTAFile interface {
	ReadExtFile() ([]byte, error)
}

// Source - file path
type Source struct {
	Path     string
	Filename string
}

// Parse MTA file and provide mta object with data
func (mta *MTA) Parse(yamlContent []byte) (err error) {
	// Format the YAML to struct's
	err = yaml.Unmarshal([]byte(yamlContent), &mta)
	if err != nil {
		return fmt.Errorf("not able to read the mta file : %s", err.Error())
	}
	return nil
}

// Marshal - usage for edit purpose
func Marshal(in MTA) (mtads []byte, err error) {
	mtads, err = yaml.Marshal(&in)
	if err != nil {
		return nil, err
	}
	return mtads, nil
}

// ReadExtFile - read external
func (s Source) ReadExtFile() ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Read MTA file
	yamlFile, err := ioutil.ReadFile(wd + pathSep + s.Path + pathSep + s.Filename)
	if err != nil {
		return yamlFile, fmt.Errorf("not able to read the mta file : %s", err.Error())
	}
	return yamlFile, err
}

// GetModules - Get list of mta modules
func (mta *MTA) GetModules() []*Modules {
	return mta.Modules
}

// GetResources - Get list of mta resources
func (mta *MTA) GetResources() []*Resources {
	return mta.Resources
}

// GetModuleByName - Get specific module
func (mta *MTA) GetModuleByName(name string) (*Modules, error) {
	for _, m := range mta.Modules {
		if m.Name == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

// GetResourceByName - Get specific resource
func (mta *MTA) GetResourceByName(name string) (*Resources, error) {
	for _, r := range mta.Resources {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

func modules(mta *MTA) []string {
	var mNames []string
	for _, mod := range mta.Modules {
		mNames = append(mNames, mod.Name)
	}
	return mNames
}

// GetModulesNames - get list of modules names
func (mta *MTA) GetModulesNames() []string {
	return modules(mta)
}

func Validate(yamlContent []byte, projectPath string, validateSchema bool, validateProject bool) mta_validate.YamlValidationIssues {
	issues := []mta_validate.YamlValidationIssue{}
	wd, err := os.Getwd()
	if err != nil {
		return append(issues, []mta_validate.YamlValidationIssue{{"Validation failed" + err.Error()}}...)
	}
	if validateSchema {
		schemaContent, _ := ioutil.ReadFile(filepath.Join(wd, "schema.yaml"))
		validations, schemaValidationLog := mta_validate.BuildValidationsFromSchemaText(schemaContent)
		if len(schemaValidationLog) > 0 {
			return schemaValidationLog
		} else {
			yamlValidationLog, err := mta_validate.ValidateYaml(yamlContent, validations...)
			if err != nil && len(yamlValidationLog) == 0 {
				yamlValidationLog = append(yamlValidationLog, []mta_validate.YamlValidationIssue{{"Validation failed" + err.Error()}}...)
			}
			issues = append(issues, yamlValidationLog...)
		}
	}
	if validateProject {
		mta := MTA{}
		yaml.Unmarshal(yamlContent, &mta)
		projectIssues := ValidateYamlProject(mta, projectPath)
		issues = append(issues, projectIssues...)
	}

	return issues
}
