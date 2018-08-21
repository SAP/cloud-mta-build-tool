package builders

//Builders list of commands to execute
type Builders struct {
	Version  string    `yaml:"version"`
	Builders []builder `yaml:"builders"`
}

type builder struct {
	Name string     `yaml:"name"`
	Type []Commands `yaml:"type"`
	Info string     `yaml:"info"`
}

type Commands struct {
	Command string `yaml:"command"`
}
