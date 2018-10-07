package mta

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"cloud-mta-build-tool/cmd/fsys"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func Test_ValidateYamlProject(t *testing.T) {
	mtaContent, _ := ioutil.ReadFile(filepath.Join(dir.GetPath(), "testdata", "testproject", "mta.yaml"))
	mta := MTA{}
	yaml.Unmarshal(mtaContent, &mta)
	issues := ValidateYamlProject(mta, filepath.Join(dir.GetPath(), "testdata", "testproject"))
	assert.Equal(t, issues[0].Msg, "Module <ui5app2> not found in project. Expected path: <ui5app2>")

}
