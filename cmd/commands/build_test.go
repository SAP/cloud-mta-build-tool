package commands

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
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

	m := mta.MTA{}
	// parse mta yaml
	err := yaml.Unmarshal(mtaCF, &m)
	if err != nil {
		logs.Logger.Error(err)
	}

	type args struct {
		mta        mta.MTA
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
				mta        mta.MTA
				moduleName string
			}{mta: m, moduleName: "htmlapp"},
			expected: []string{"npm install", "grunt", "npm prune --production"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := moduleCmd(tt.args.mta, tt.args.moduleName); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("moduleCmd() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_cmdConverter(t *testing.T) {

	cmdParams := [][]string{
		{"path", "npm", "install"},
		{"path", "grunt"},
		{"path", "npm", "prune", "--production"},
	}

	type args struct {
		mPath   string
		cmdList []string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Command converter",
			args: struct {
				mPath   string
				cmdList []string
			}{mPath: "path", cmdList: []string{"npm install", "grunt", "npm prune --production"}},
			want: cmdParams,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmdConverter(tt.args.mPath, tt.args.cmdList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdConverter() = %v,`\n\t` want %v", got, tt.want)
			}
		})
	}
}
