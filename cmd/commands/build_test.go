package commands

import (
	"reflect"
	"testing"

	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta/models"

	"gopkg.in/yaml.v2"
)

func Test_moduleCmd(t *testing.T) {

	var mtaCF = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app

  - name: htmlapp2
    type: html5
    path: app

  - name: java
    type: java
    path: app
`)

	mta := models.MTA{}
	// parse mta yaml
	err := yaml.Unmarshal(mtaCF, &mta)
	if err != nil {
		logs.Logger.Error(err)
	}

	type args struct {
		mta        models.MTA
		moduleName string
	}
	tests := []struct {
		name     string
		args     args
		expected []string
	}{
		{
			name: "build specific module by name",
			args: struct {
				mta        models.MTA
				moduleName string
			}{mta: mta, moduleName: "htmlapp"},
			expected: []string{"npm install", "grunt", "npm prune --production"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := moduleCmd(tt.args.mta, tt.args.moduleName); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("moduleCmd() = %v, want %v", got, tt.expected)
			}
		})
	}
}
