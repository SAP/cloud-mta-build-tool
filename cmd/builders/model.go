package builders

//Builders list of commands to execute
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

//Commands - specific command
type Commands struct {
	Command string `yaml:"command"`
}
