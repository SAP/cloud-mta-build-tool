package builders

import (
	"log"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/mta"
)

func TestExeCmd(t *testing.T) {
	t.Parallel()
	var buildersCfg = []byte(`
version: 1
builders:
  - name: html5
    info: "build UI application"
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

	// Get parsed yaml content
	commands := Builders{}
	err := yaml.Unmarshal(buildersCfg, &commands)
	if err != nil {
		log.Println("Error: " + err.Error())
	}

	type args struct {
		modules mta.Modules
	}
	tests := []struct {
		name     string
		args     args
		expected commandList
	}{
		{
			name: "Command for required module type",
			args: args{
				modules: mta.Modules{
					Name:       "uiapp",
					Type:       "html5",
					Path:       "./",
					Requires:   nil,
					Provides:   nil,
					Parameters: nil,
					Properties: nil,
				},
			},
			expected: commandList{
				Info:    "build UI application",
				Command: []string{"npm install", "grunt", "npm prune --production"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mesh(tt.args.modules, commands); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("CommandProvider() = \n got \t\n %v, \n expected \t\n %v", got, tt.expected)
			}
		})
	}
}

func TestCommandProvider(t *testing.T) {
	type args struct {
		modules mta.Modules
	}
	tests := []struct {
		name string
		args args
		want commandList
	}{
		{
			name: "Command for required module type",
			args: args{
				modules: mta.Modules{
					Type: "html5",
				},
			},
			want: struct {
				Info    string
				Command []string
			}{
				Info:    "installing module dependencies & execute grunt & remove dev dependencies",
				Command: []string{"npm install", "grunt", "npm prune --production"},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := CommandProvider(tt.args.modules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
