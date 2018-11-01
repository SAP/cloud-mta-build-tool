package version

import (
	"gopkg.in/yaml.v2"
)

type Version struct {
	CliVersion string `yaml:"cli_version"`
	MakeFile   string `yaml:"makefile_version"`
}

func GetVersion() (Version, error) {
	v := Version{}
	err := yaml.Unmarshal(VersionConfig, &v)
	return v, err
}
