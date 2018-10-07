package mta

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/validations"
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
type Properties map[string]string

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

type file interface {
	ReadExtFile() ([]byte, error)
}

// Source - file path
type Source struct {
	Path     string
	Filename string
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

func Validate(yamlContent []byte) bool {
	schemaContent, _ := ioutil.ReadFile(filepath.Join(dir.GetPath(), "schema.yaml"))
	validations, schemaValidationLog := mta_validate.BuildValidationsFromSchemaText(schemaContent)
	if len(schemaValidationLog) > 0 {
		logs.Logger.Error(schemaValidationLog)
		return false
	} else {
		issues, err := mta_validate.ValidateYaml(yamlContent, validations...)
		if err != nil || len(issues) > 0 {
			logs.Logger.Error(issues)
			return false
		}
		return true
	}
}

// ReadExtFile - read external
func (s Source) ReadExtFile() ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		logs.Logger.Error(err)
	}
	// Read MTA file
	yamlFile, err := ioutil.ReadFile(wd + pathSep + s.Path + pathSep + s.Filename)
	if err != nil {
		return yamlFile, fmt.Errorf("not able to read the mta file : %s", err.Error())
	}
	return yamlFile, err
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
