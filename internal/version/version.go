package version

import (
	"fmt"

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

// GetVersionMessage returns the message for the "version" flag
func GetVersionMessage() (string, error) {
	v, err := GetVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(versionMsg, v.CliVersion), nil
}
