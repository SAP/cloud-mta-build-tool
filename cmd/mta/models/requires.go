package models

// Requires / Mta struct
type Requires struct {
	Name       string     `yaml:"name,omitempty"`
	Group      string     `yaml:"group,omitempty"`
	Type       string     `yaml:"type,omitempty"`
	Properties Properties `yaml:"properties,omitempty"`
}

type BuildRequires struct {
	Name       string     `yaml:"name,omitempty"`
	TargetPath string     `yaml:"target-path,omitempty"`
}
