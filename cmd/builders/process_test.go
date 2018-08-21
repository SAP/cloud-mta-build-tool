package builders

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/logs"
)

func TestParse(t *testing.T) {

	var buildCfg = []byte(`
version: 1
builders:
  - name: html5
    info: "build UI5 application"
    type:
    - command: npm install
    - command: grunt
    - command: npm prune --production
  - name: java
    info: "build java application"
    type:
    - command: mvn clean install
  - name: nodejs
    info: "build nodejs application"
    type:
    - command: npm install
  - name: golang
    info: "build golang application"
    type:
    - command: go build *.go
`)

	var wantOut = []byte(`
version: 1
builders:
  - name: html5
    info: "build UI5 application"
    type:
    - command: npm install
    - command: grunt
    - command: npm prune --production
  - name: java
    info: "build java application"
    type:
    - command: mvn clean install
  - name: nodejs
    info: "build nodejs application"
    type:
    - command: npm install
  - name: golang
    info: "build golang application"
    type:
    - command: go build *.go
`)

	//Get parsed yaml content
	commands := Builders{}
	err := yaml.Unmarshal(wantOut, &commands)
	if err != nil {
		logs.Logger.Error("Error: " + err.Error())
	}

	tests := []struct {
		name     string
		args     []byte
		expected Builders
	}{
		{
			name:     "Parse builders configuration files",
			args:     buildCfg,
			expected: commands,
		},
	}
	//Todo - basic parse test, need types test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() = %v, \n expected %v", got, tt.expected)
			}
		})
	}
}
