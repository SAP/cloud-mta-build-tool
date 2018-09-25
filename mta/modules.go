package mta

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
