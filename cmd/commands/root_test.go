package commands

import (
	"path/filepath"
	"testing"

	"cloud-mta-build-tool/cmd/fsys"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_initConfig(t *testing.T) {

	initConfig()
	property := viper.Get("xxx")
	assert.Nil(t, property)

	cfgFile = filepath.Join(dir.GetPath(), "testdata", "config.props")
	initConfig()
	property = viper.Get("xxx")
	assert.Equal(t, "10", property)

	viper.Reset()

	cfgFile = filepath.Join(dir.GetPath(), "testdata", "config1.props")
	initConfig()
	cfgFile = ""
	property = viper.Get("xxx")
	assert.Nil(t, property)

}
