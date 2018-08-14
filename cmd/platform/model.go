package platform

//Platforms - data structure for platforms module types configuration
type Platforms struct {
	Platforms []Modules `yaml:"platform"`
}
type Modules struct {
	Name   string       `yaml:"name"`
	Models []Properties `yaml:"modules"`
}

type Properties struct {
	NativeType   string `yaml:"native-type"`
	PlatformType string `yaml:"platform-type"`
}
