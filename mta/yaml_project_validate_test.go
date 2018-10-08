package mta

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func Test_ValidateYamlProject(t *testing.T) {

	wd, _ := os.Getwd()
	mtaContent, _ := ioutil.ReadFile(filepath.Join(wd, "testdata", "testproject", "mta.yaml"))
	mta := MTA{}
	yaml.Unmarshal(mtaContent, &mta)
	issues := ValidateYamlProject(mta, filepath.Join(wd, "testdata", "testproject"))
	assert.Equal(t, issues[0].Msg, "Module <ui5app2> not found in project. Expected path: <ui5app2>")
}
