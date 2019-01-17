package commands

// Builders list of commands to execute
type Builders struct {
	Version  string    `yaml:"version"`
	Builders []builder `yaml:"builders"`
}

type builder struct {
	Name string     `yaml:"name"`
	Info string     `yaml:"info"`
	Path string     `yaml:"path"`
	Type []Commands `yaml:"type"`
}

// Commands - specific command
type Commands struct {
	Command string `yaml:"command"`
}

// ModuleTypes - list of commands or builders related to specific module type
type ModuleTypes struct {
	Version     string       `yaml:"version"`
	ModuleTypes []moduleType `yaml:"module-types"`
}

type moduleType struct {
	Name    string     `yaml:"name"`
	Info    string     `yaml:"info"`
	Path    string     `yaml:"path"`
	Builder string     `yaml:"builder"`
	Type    []Commands `yaml:"type"`
}
