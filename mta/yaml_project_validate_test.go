package mta

import (
	"io/ioutil"
	"testing"

	"cloud-mta-build-tool/cmd/fsys"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func Test_ValidateYamlProject(t *testing.T) {
	mtaYamlPath, _ := dir.GetFullPath("testdata", "testproject", "mta.yaml")
	mtaContent, _ := ioutil.ReadFile(mtaYamlPath)
	mta := MTA{}
	yaml.Unmarshal(mtaContent, &mta)
	mtaProjectPath, _ := dir.GetFullPath("testdata", "testproject")
	issues := ValidateYamlProject(mta, mtaProjectPath)
	assert.Equal(t, issues[0].Msg, "Module <ui5app2> not found in project. Expected path: <ui5app2>")
}
