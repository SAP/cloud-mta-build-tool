package builders

//ExeCommands list of commands to execute
type ExeCommands struct {
	Version  string     `yaml:"version"`
	Builders []builders `yaml:"builders"`
}

type builders struct {
	Name string    `yaml:"name"`
	Type []command `yaml:"type"`
	Info string    `yaml:"info"`
}

type command struct {
	Command string `yaml:"command"`
}
