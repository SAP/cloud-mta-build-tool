package platform

//Platforms - data structure for platforms module types configuration
type Platforms struct {
	Platforms map[string]modules `yaml:"platforms"`
}
type modules map[string][]properties

type properties struct {
	NativeType   string `yaml:"native-type"`
	PlatformType string `yaml:"platform-type"`
}
