package models

type Resources struct {
	Name       string
	Type       string
	Parameters Parameters `yaml:"parameters,omitempty"`
	Properties Properties `yaml:"properties,omitempty"`
}
