package mta

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mbtv2/cmd/mta/models"
)
// Table driven test
// Unit test for parsing mta files to working object
func Test_ParseFile(t *testing.T) {



	tests := []struct {
		n       int
		name    string
		wantOut models.Modules
	}{
		{
			name: "Parse service(srv) Module section",
			wantOut: models.Modules{
				Name: "srv",
				Type: "java",
				Path: "srv",
				Requires: []models.Requires{
					{
						Name: "db",
						Properties: models.Properties{
							"JBP_CONFIG_RESOURCE_CONFIGURATION": `[tomcat/webapps/ROOT/META-INF/context.xml: {"service_name_for_DefaultDB" : "~{hdi-container-name}"}]`,
						},
					},
				},
				Provides: []models.Provides{
					{

						Name: "srv_api",
						Properties: models.Properties{
							"url": "${default-url}",
						},
					},
				},
				Parameters: models.Parameters{
					"memory": "512M",
				},
				// BuildParams: nil,
				Properties: models.Properties{
					"APPC_LOG_LEVEL":              "info",
					"VSCODE_JAVA_DEBUG_LOG_LEVEL": "ALL",
				},
			},
		},

		// ------------------------Second module test------------------------------
		{
			name: "Parse UI(HTML5) Module section",
			wantOut: models.Modules{
				Name: "ui",
				Type: "html5",
				Path: "ui",
				Requires: []models.Requires{
					{
						Name:  "srv_api",
						Group: "destinations",
						Properties: models.Properties{
							"forwardAuthToken": "true",
							"strictSSL":        "false",
							"name":             "srv_api",
							"url":              "~{url}",
						},
					},
				},
				Provides: []models.Provides{
					{

						Name: "srv_api",
						Properties: models.Properties{
							"url": "${default-url}",
						},
					},
				},
				BuildParams: models.BuildParameters{

					"builder": "grunt",
				},

				Parameters: models.Parameters{
					"disk-quota": "256M",
					"memory":     "256M",
				},
			},
		},
	}

	// First Module test as atomic building blocks

	mtaFile, _ := ioutil.ReadFile("./testdata/mta.yaml")
	var idx int
	actual, err := Parse(mtaFile)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Switch was added to handle different type of slices
			switch idx {
			// Run Service module
			case 0:
				require.NoError(t, err)
				require.NotNil(t, actual)
				require.Len(t, actual.Modules, 2)
				assert.Equal(t, tt.wantOut.Name, actual.Modules[tt.n].Name)
				assert.Equal(t, tt.wantOut.Type, actual.Modules[tt.n].Type)
				assert.Equal(t, tt.wantOut.Path, actual.Modules[tt.n].Path)
				assert.Equal(t, tt.wantOut.Parameters, actual.Modules[tt.n].Parameters)
				assert.Equal(t, tt.wantOut.Properties, actual.Modules[tt.n].Properties)
				assert.Equal(t, tt.wantOut.Requires, actual.Modules[tt.n].Requires)
				assert.Equal(t, tt.wantOut.Provides, actual.Modules[tt.n].Provides)

				// Run UI module
			case 1:

				assert.Equal(t, tt.wantOut.Name, actual.Modules[idx].Name)
				assert.Equal(t, tt.wantOut.Type, actual.Modules[idx].Type)
				assert.Equal(t, tt.wantOut.Path, actual.Modules[idx].Path)
				assert.Equal(t, tt.wantOut.Requires, actual.Modules[idx].Requires)
				assert.Equal(t, tt.wantOut.Parameters, actual.Modules[idx].Parameters)
				assert.Equal(t, tt.wantOut.BuildParams, actual.Modules[idx].BuildParams)

			}

			idx = idx + 1

		})

	}

}
