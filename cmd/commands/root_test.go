package commands

import (
	"testing"

	"cloud-mta-build-tool/internal/fsys"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_initConfig(t *testing.T) {

	initConfig()
	property := viper.Get("xxx")
	assert.Nil(t, property)

	cfgFile, _ = dir.GetFullPath("testdata", "config.props")
	initConfig()
	property = viper.Get("xxx")
	assert.Equal(t, "10", property)

	viper.Reset()

	cfgFile, _ = dir.GetFullPath("testdata", "config1.props")
	initConfig()
	cfgFile = ""
	property = viper.Get("xxx")
	assert.Nil(t, property)

}
