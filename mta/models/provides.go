package models

// Provides - MTA struct
type Provides struct {
	Name       string
	Properties Properties `yaml:"properties,omitempty"`
}
