package builders

import (
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta/models"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestExeCmd(t *testing.T) {

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

	//Get parsed yaml content
	commands := Builders{}
	err := yaml.Unmarshal(buildersCfg, &commands)
	if err != nil {
		logs.Logger.Error("Error: " + err.Error())
	}

	type args struct {
		modules models.Modules
	}
	tests := []struct {
		name string
		args args
		want []CommandList
	}{
		{
			name: "Command for required module type",
			args: args{
				modules: models.Modules{
					Name:       "uiapp",
					Type:       "html5",
					Path:       "./",
					Requires:   nil,
					Provides:   nil,
					Parameters: nil,
					Properties: nil,
				},
			},
			want: []CommandList{
				{"build UI application",
					[]string{"npm install", "grunt", "npm prune --production"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mesh(tt.args.modules, commands); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandProvider() = \n got \t\n %v, \n want \t\n %v", got, tt.want)
			}
		})
	}
}

func TestCommandProvider(t *testing.T) {
	type args struct {
		modules models.Modules
	}
	tests := []struct {
		name string
		args args
		want []CommandList
	}{
		{
			name: "Command for required module type",
			args: args{
				modules: models.Modules{
					Type: "html5",
				},
			},
			want: []CommandList{
				{
					"installing module dependencies & execute grunt & remove dev dependencies",
					[]string{"npm install", "grunt", "npm prune --production"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommandProvider(tt.args.modules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
