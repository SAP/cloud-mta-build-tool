package version

import (
	"gopkg.in/yaml.v2"
)

// Version - tool version
type Version struct {
	CliVersion string `yaml:"cli_version"`
	MakeFile   string `yaml:"makefile_version"`
}

// GetVersion - get versions
func GetVersion() (Version, error) {
	v := Version{}
	err := yaml.UnmarshalStrict(VersionConfig, &v)
	return v, err
}
