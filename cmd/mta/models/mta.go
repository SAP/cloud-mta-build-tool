package models

// MTA - Main mta struct
type MTA struct {
	SchemaVersion *string      `yaml:"_schema-version"`
	Id            string       `yaml:"ID"`
	Version       string       `yaml:"version,omitempty"`
	Modules       []*Modules   `yaml:"modules,omitempty"`
	Resources     []*Resources `yaml:"resources,omitempty"`
}
