package mta

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/validations"
)

// MTA struct
type MTA struct {
	// Todo 1. Provide interface to support multiple mta schema (2.1 / 3.1 ) versions and concrete struct type
	// Todo 2. Add missing properties
	// indicates MTA schema version, using semver
	SchemaVersion *string `yaml:"_schema-version"`
	// A globally unique ID of this MTA. Unlimited string of unicode characters.
	ID string `yaml:"ID"`
	// A non-translatable description of this MTA. This is not a text for application users
	Description string `yaml:"description,omitempty"`
	// Application version, using semantic versioning standard
	Version string `yaml:"version,omitempty"`
	// The provider or vendor of this software
	Provider string `yaml:"provider,omitempty"`
	// A copyright statement from the provider
	Copyright string `yaml:"copyright,omitempty"`
	// list of modules
	Modules []*Modules `yaml:"modules,omitempty"`
	// Resource declarations. Resources can be anything required to run the application which is not provided by the application itself
	Resources []*Resources `yaml:"resources,omitempty"`
	// Parameters can be used to steer the behavior of tools which interpret this descriptor
	Parameters Parameters `yaml:"parameters,omitempty"`
}

// Build-parameters are specifically steering the behavior of build tools.
type buildParameters struct {
	// Builder name
	Builder string `yaml:"builder,omitempty"`
	// Builder type
	Type string `yaml:"type,omitempty"`
	// A path pointing to a file which contains a map of parameters, either in JSON or in YAML format.
	Path string `yaml:"path,omitempty"`
	// list of names either matching a resource name or a name provided by another module within the same MTA
	Requires []BuildRequires `yaml:"requires,omitempty"`
}

// Modules - MTA modules
type Modules struct {
	// An MTA internal module name. Names need to be unique within the MTA scope
	Name string
	// a globally unique type ID. Deployment tools will interpret this type ID
	Type string
	// A file path which identifies the location of module artifacts.
	Path string `yaml:"path,omitempty"`
	// list of names either matching a resource name or a name provided by another module within the same MTA
	Requires []Requires `yaml:"requires,omitempty"`
	// List of provided names (MTA internal)to which properties (= configuration data) can be attached
	Provides []Provides `yaml:"provides,omitempty"`
	// Parameters can be used to steer the behavior of tools which interpret this descriptor. Parameters are not made available to the module at runtime
	Parameters Parameters `yaml:"parameters,omitempty"`
	// Build-parameters are specifically steering the behavior of build tools.
	BuildParams buildParameters `yaml:"build-parameters,omitempty"`
	// Provided property values can be accessed by "~{<name-of-provides-section>/<provided-property-name>}". Such expressions can be part of an arbitrary string
	Properties Properties `yaml:"properties,omitempty"`
}

// Properties - properties map
type Properties map[string]interface{}

// Parameters - parameters map
type Parameters map[string]interface{}

// Provides List of provided names (MTA internal)to which properties (= configuration data) can be attached
type Provides struct {
	Name       string
	Properties Properties `yaml:"properties,omitempty"`
}

// Requires List of names either matching a resource name or a name provided by another module within the same MTA
type Requires struct {
	// an MTA internal name which must match either a provided name, a resource name, or a module name within the same MTA
	Name string `yaml:"name,omitempty"`
	// A group name which shall be use by a deployer to group properties for lookup by a module runtime.
	Group string `yaml:"group,omitempty"`
	Type  string `yaml:"type,omitempty"`
	// Provided property values can be accessed by "~{<provided-property-name>}". Such expressions can be part of an arbitrary string
	Properties Properties `yaml:"properties,omitempty"`
}

// BuildRequires - build requires section
type BuildRequires struct {
	Name       string   `yaml:"name,omitempty"`
	Artifacts  []string `yaml:"artifacts,omitempty"`
	TargetPath string   `yaml:"target-path,omitempty"`
}

// Resources can be anything required to run the application which is not provided by the application itself.
type Resources struct {
	Name string
	// A type of a resource. This type is interpreted by and must be known to the deployer. Resources can be untyped
	Type string
	// Parameters can be used to influence the behavior of tools which interpret this descriptor. Parameters are not made available to requiring modules at runtime
	Parameters Parameters `yaml:"parameters,omitempty"`
	// property names and values make up the configuration data which is to be provided to requiring modules at runtime
	Properties Properties `yaml:"properties,omitempty"`
}

// Parse parses MTA YAML document and provides MTA object with data
func (mta *MTA) Parse(yamlContent []byte) (err error) {
	// Format the YAML to struct's
	err = yaml.Unmarshal([]byte(yamlContent), &mta)
	if err != nil {
		return errors.Wrap(err, "error occurred while parsing file : %s")
	}
	return nil
}

// Marshal serializes the MTA provided into a YAML document, edit scenario
func Marshal(in *MTA) (mtads []byte, err error) {
	mtads, err = yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	return mtads, nil
}

// ReadMtaYaml Reads MTA Yaml file and stores it's data in byte slice.
func ReadMtaYaml(ep *fs.MtaLocationParameters) ([]byte, error) {
	fileFullPath, err := ep.GetMtaYamlPath()
	if err != nil {
		return nil, errors.Wrap(err, "ReadMtaYaml failed getting MTA Yaml path")
	}
	// Read MTA file
	yamlFile, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		return nil, errors.Wrap(err, "ReadMtaYaml failed getting MTA Yaml path reading the mta file")
	}
	return yamlFile, nil
}

// GetModules Returns list of mta modules
func (mta *MTA) GetModules() []*Modules {
	return mta.Modules
}

// GetResources Returns list of mta resources
func (mta *MTA) GetResources() []*Resources {
	return mta.Resources
}

// GetModuleByName Returns specific module
func (mta *MTA) GetModuleByName(name string) (*Modules, error) {
	for _, m := range mta.Modules {
		if m.Name == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

// GetResourceByName Returns specific resource
func (mta *MTA) GetResourceByName(name string) (*Resources, error) {
	for _, r := range mta.Resources {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("module %s , not found ", name)
}

// GetModulesNames Returns list of modules names
func (mta *MTA) GetModulesNames() ([]string, error) {
	return mta.getModulesOrder()
}

// Validate Validate mta schema
func Validate(yamlContent []byte, projectPath string, validateSchema bool, validateProject bool) validate.YamlValidationIssues {
	//noinspection GoPreferNilSlice
	issues := []validate.YamlValidationIssue{}
	if validateSchema {
		validations, schemaValidationLog := validate.BuildValidationsFromSchemaText(schemaDef)
		if len(schemaValidationLog) > 0 {
			return schemaValidationLog
		} else {
			yamlValidationLog, err := validate.ValidateYaml(yamlContent, validations...)
			if err != nil && len(yamlValidationLog) == 0 {
				yamlValidationLog = append(yamlValidationLog, []validate.YamlValidationIssue{{Msg: "Validation failed" + err.Error()}}...)
			}
			issues = append(issues, yamlValidationLog...)
		}
	}
	if validateProject {
		mta := MTA{}
		yaml.Unmarshal(yamlContent, &mta)
		projectIssues := validateYamlProject(&mta, projectPath)
		issues = append(issues, projectIssues...)
	}

	return issues
}
