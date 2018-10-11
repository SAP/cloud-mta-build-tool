package mta

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testInfo struct {
	name      string
	expected  Modules
	validator func(t *testing.T, actual, expected Modules)
}

func doTest(t *testing.T, expected []testInfo, filename string) {

	mtaFile, _ := ioutil.ReadFile(filename)
	// Parse file
	oMta := &MTA{}
	err := oMta.Parse(mtaFile)
	for i, tt := range expected {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, oMta)
			require.Len(t, oMta.Modules, len(expected))
			tt.validator(t, *oMta.Modules[i], tt.expected)
		})
	}

	mtaContent, err := Marshal(*oMta)
	assert.Nil(t, err)

	oMta2 := &MTA{}
	newErr := oMta2.Parse(mtaContent)
	assert.Nil(t, newErr)
	assert.Equal(t, oMta, oMta2)
}

func Test_ValidateAll(t *testing.T) {

	wd, _ := os.Getwd()
	yamlContent, _ := ioutil.ReadFile(filepath.Join(wd, "testdata", "testproject", "mta.yaml"))
	issues := Validate(yamlContent, (filepath.Join(wd, "testdata", "testproject")), true, true)
	assert.Equal(t, 1, len(issues))
}

func Test_ValidateSchema(t *testing.T) {
	wd, _ := os.Getwd()
	yamlContent, _ := ioutil.ReadFile(filepath.Join(wd, "testdata", "mta_multiapps.yaml"))
	issues := Validate(yamlContent, (filepath.Join(wd, "testdata")), true, false)
	assert.Equal(t, 0, len(issues))
}

// Table driven test
// Unit test for parsing mta files to working object
func Test_ModulesParsing(t *testing.T) {
	tests := []testInfo{
		{
			name: "Parse service(srv) Module section",
			expected: Modules{
				Name: "srv",
				Type: "java",
				Path: "srv",
				Requires: []Requires{
					{
						Name: "db",
						Properties: Properties{
							"JBP_CONFIG_RESOURCE_CONFIGURATION": `[tomcat/webapps/ROOT/META-INF/context.xml: {"service_name_for_DefaultDB" : "~{hdi-container-name}"}]`,
						},
					},
				},
				Provides: []Provides{
					{
						Name: "srv_api",
						Properties: Properties{
							"url": "${default-url}",
						},
					},
				},
				Parameters: Parameters{
					"memory": "512M",
				},
				Properties: Properties{
					"APPC_LOG_LEVEL":              "info",
					"VSCODE_JAVA_DEBUG_LOG_LEVEL": "ALL",
				},
			},
			validator: func(t *testing.T, actual, expected Modules) {
				assert.Equal(t, expected.Name, actual.Name)
				assert.Equal(t, expected.Type, actual.Type)
				assert.Equal(t, expected.Path, actual.Path)
				assert.Equal(t, expected.Parameters, actual.Parameters)
				assert.Equal(t, expected.Properties, actual.Properties)
				assert.Equal(t, expected.Requires, actual.Requires)
				assert.Equal(t, expected.Provides, actual.Provides)
			}},

		// ------------------------Second module test------------------------------
		{
			name: "Parse UI(HTML5) Module section",
			expected: Modules{
				Name: "ui",
				Type: "html5",
				Path: "ui",
				Requires: []Requires{
					{
						Name:  "srv_api",
						Group: "destinations",
						Properties: Properties{
							"forwardAuthToken": true,
							"strictSSL":        false,
							"name":             "srv_api",
							"url":              "~{url}",
						},
					},
				},
				Provides: []Provides{
					{
						Name: "srv_api",
						Properties: Properties{
							"url": "${default-url}",
						},
					},
				},
				BuildParams: BuildParameters{
					Builder: "grunt",
				},

				Parameters: Parameters{
					"disk-quota": "256M",
					"memory":     "256M",
				},
			},
			validator: func(t *testing.T, actual, expected Modules) {
				assert.Equal(t, expected.Name, actual.Name)
				assert.Equal(t, expected.Type, actual.Type)
				assert.Equal(t, expected.Path, actual.Path)
				assert.Equal(t, expected.Requires, actual.Requires)
				assert.Equal(t, expected.Parameters, actual.Parameters)
				assert.Equal(t, expected.BuildParams, actual.BuildParams)
			}},
	}

	doTest(t, tests, "./testdata/mta.yaml")

}

func Test_BrokenMta(t *testing.T) {

	mtaContent, _ := ioutil.ReadFile("./testdata/mtaWithBrokenProperties.yaml")

	oMta := &MTA{}
	err := oMta.Parse(mtaContent)
	require.NotNil(t, err)
	require.NotNil(t, oMta)
}

func Test_FullMta(t *testing.T) {
	schemaVersion := "2.0.0"

	expected := MTA{
		SchemaVersion: &schemaVersion,
		Id:            "cloud.samples.someproj",
		Version:       "1.0.0",
		Parameters: Parameters{
			"deploy_mode": "html5-repo",
		},
		Modules: []*Modules{
			{
				Name: "someproj-db",
				Type: "hdb",
				Path: "db",
				Requires: []Requires{
					{
						Name: "someproj-hdi-container",
					},
					{
						Name: "someproj-logging",
					},
				},
				Parameters: Parameters{
					"disk-quota": "256M",
					"memory":     "256M",
				},
			},
			{
				Name: "someproj-java",
				Type: "java",
				Path: "srv",
				Parameters: Parameters{
					"memory":     "512M",
					"disk-quota": "256M",
				},
				Provides: []Provides{
					{
						Name: "java",
						Properties: Properties{
							"url": "${default-url}",
						},
					},
				},
				Requires: []Requires{
					{
						Name: "someproj-hdi-container",
						Properties: Properties{
							"JBP_CONFIG_RESOURCE_CONFIGURATION": "[tomcat/webapps/ROOT/META-INF/context.xml: " +
								"{\"service_name_for_DefaultDB\" : \"~{hdi-container-name}\"}]",
						},
					},
					{
						Name: "someproj-logging",
					},
				},
				BuildParams: BuildParameters{
					Requires: []BuildRequires{
						{
							Name:       "someproj-db",
							TargetPath: "",
						},
					},
				},
			},
			{
				Name: "someproj-catalog-ui",
				Type: "html5",
				Path: "someproj-someprojCatalog",
				Parameters: Parameters{
					"memory":     "256M",
					"disk-quota": "256M",
				},
				Requires: []Requires{
					{
						Name:  "java",
						Group: "destinations",
						Properties: Properties{
							"name": "someproj-backend",
							"url":  "~{url}",
						},
					},
					{
						Name: "someproj-logging",
					},
				},
				BuildParams: BuildParameters{
					Builder: "grunt",
					Requires: []BuildRequires{
						{
							Name:       "someproj-java",
							TargetPath: "",
						},
					},
				},
			},
			{
				Name: "someproj-uideployer",
				Type: "com.sap.html5.application-content",
				Parameters: Parameters{
					"memory":     "256M",
					"disk-quota": "256M",
				},
				Requires: []Requires{
					{
						Name: "someproj-apprepo-dt",
					},
				},
				BuildParams: BuildParameters{
					Builder: "grunt",
					Type:    "com.sap.html5.application-content",
					Requires: []BuildRequires{
						{
							Name: "someproj-catalog-ui",
						},
					},
				},
			},
			{
				Name: "someproj",
				Type: "approuter.nodejs",
				Path: "approuter",
				Parameters: Parameters{
					"memory":     "256M",
					"disk-quota": "256M",
				},
				Requires: []Requires{
					{
						Name:  "java",
						Group: "destinations",
						Properties: Properties{
							"name": "someproj-backend",
							"url":  "~{url}",
						},
					},
					{
						Name: "someproj-apprepo-rt",
					},
					{
						Name: "someproj-logging",
					},
				},
			},
		},
		Resources: []*Resources{
			{
				Name: "someproj-hdi-container",
				Properties: Properties{
					"hdi-container-name": "${service-name}",
				},
				Type: "com.sap.xs.hdi-container",
			},
			{
				Name: "someproj-apprepo-rt",
				Type: "org.cloudfoundry.managed-service",
				Parameters: Parameters{
					"service":      "html5-apps-repo",
					"service-plan": "app-runtime",
				},
			},
			{
				Name: "someproj-apprepo-dt",
				Type: "org.cloudfoundry.managed-service",
				Parameters: Parameters{
					"service":      "html5-apps-repo",
					"service-plan": "app-host",
				},
			},
			{
				Name: "someproj-logging",
				Type: "org.cloudfoundry.managed-service",
				Parameters: Parameters{
					"service":      "application-logs",
					"service-plan": "lite",
				},
			},
		},
	}

	mtaContent, _ := ioutil.ReadFile("./testdata/mta2.yaml")

	actual := &MTA{}
	err := actual.Parse(mtaContent)
	assert.Nil(t, err)
	assert.Equal(t, expected, *actual)

}

func TestMTA_GetModules(t *testing.T) {
	type fields struct {
		SchemaVersion *string
		Id            string
		Version       string
		Modules       []*Modules
		Resources     []*Resources
		Parameters    Parameters
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Modules
	}{
		{
			name: "GetModules - two Modules",
			fields: fields{
				Modules: []*Modules{
					{
						Name: "someproj-db",
						Type: "hdb",
						Path: "db",
						Requires: []Requires{
							{
								Name: "someproj-hdi-container",
							},
							{
								Name: "someproj-logging",
							},
						},
					},
					{
						Name: "someproj-java",
						Type: "java",
						Path: "srv",
						Parameters: Parameters{
							"memory":     "512M",
							"disk-quota": "256M",
						},
					},
				},
				Resources: []*Resources{
					{
						Name: "someproj-hdi-container",
						Properties: Properties{
							"hdi-container-name": "${service-name}",
						},
						Type: "com.sap.xs.hdi-container",
					},
				},
			},
			want: []*Modules{
				{
					Name: "someproj-db",
					Type: "hdb",
					Path: "db",
					Requires: []Requires{
						{
							Name: "someproj-hdi-container",
						},
						{
							Name: "someproj-logging",
						},
					},
				},
				{
					Name: "someproj-java",
					Type: "java",
					Path: "srv",
					Parameters: Parameters{
						"memory":     "512M",
						"disk-quota": "256M",
					},
				},
			},
		}, {
			name:   "GetModules - Empty list",
			fields: fields{},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mta := &MTA{
				SchemaVersion: tt.fields.SchemaVersion,
				Id:            tt.fields.Id,
				Version:       tt.fields.Version,
				Modules:       tt.fields.Modules,
				Resources:     tt.fields.Resources,
				Parameters:    tt.fields.Parameters,
			}
			got := mta.GetModules()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MTA.GetModules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMTA_GetResources(t *testing.T) {
	type fields struct {
		SchemaVersion *string
		Id            string
		Version       string
		Modules       []*Modules
		Resources     []*Resources
		Parameters    Parameters
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Resources
	}{
		{
			name: "GetResources - two resources",
			fields: fields{
				Modules: []*Modules{
					{
						Name: "someproj-db",
						Type: "hdb",
						Path: "db",
						Requires: []Requires{
							{
								Name: "someproj-hdi-container",
							},
							{
								Name: "someproj-logging",
							},
						},
					},
					{
						Name: "someproj-java",
						Type: "java",
						Path: "srv",
						Parameters: Parameters{
							"memory":     "512M",
							"disk-quota": "256M",
						},
					},
				},
				Resources: []*Resources{
					{
						Name: "someproj-hdi-container",
						Properties: Properties{
							"hdi-container-name": "${service-name}",
						},
						Type: "com.sap.xs.hdi-container",
					},
					{
						Name: "someproj-apprepo-rt",
						Type: "org.cloudfoundry.managed-service",
						Parameters: Parameters{
							"service":      "html5-apps-repo",
							"service-plan": "app-runtime",
						},
					},
				},
			},
			want: []*Resources{
				{
					Name: "someproj-hdi-container",
					Properties: Properties{
						"hdi-container-name": "${service-name}",
					},
					Type: "com.sap.xs.hdi-container",
				},
				{
					Name: "someproj-apprepo-rt",
					Type: "org.cloudfoundry.managed-service",
					Parameters: Parameters{
						"service":      "html5-apps-repo",
						"service-plan": "app-runtime",
					},
				},
			},
		}, {
			name: "GetResources - Empty list",
			fields: fields{
				SchemaVersion: nil,
				Id:            "",
				Version:       "",
				Modules:       nil,
				Resources:     nil,
				Parameters:    nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mta := &MTA{
				SchemaVersion: tt.fields.SchemaVersion,
				Id:            tt.fields.Id,
				Version:       tt.fields.Version,
				Modules:       tt.fields.Modules,
				Resources:     tt.fields.Resources,
				Parameters:    tt.fields.Parameters,
			}

			got := mta.GetResources()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MTA.GetResources() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMTA_GetModuleByName(t *testing.T) {

	type fields struct {
		SchemaVersion *string
		Id            string
		Version       string
		Modules       []*Modules
		Resources     []*Resources
		Parameters    Parameters
	}
	{
		tests := []struct {
			name       string
			fields     fields
			moduleName string
			want       *Modules
			wantErr    bool
		}{
			{
				name:       "GetModuleByName",
				moduleName: "someproj-java",
				fields: fields{
					Modules: []*Modules{
						{
							Name: "someproj-db",
							Type: "hdb",
							Path: "db",
							Requires: []Requires{
								{
									Name: "someproj-hdi-container",
								},
								{
									Name: "someproj-logging",
								},
							},
							Parameters: Parameters{
								"disk-quota": "256M",
								"memory":     "256M",
							},
						},
						{
							Name: "someproj-java",
							Type: "java",
							Path: "srv",
							Parameters: Parameters{
								"memory":     "512M",
								"disk-quota": "256M",
							},
							Provides: []Provides{
								{
									Name: "java",
									Properties: Properties{
										"url": "${default-url}",
									},
								},
							},
						}},
					Resources: nil,
				},
				want: &Modules{

					Name: "someproj-java",
					Type: "java",
					Path: "srv",
					Parameters: Parameters{
						"memory":     "512M",
						"disk-quota": "256M",
					},
					Provides: []Provides{
						{
							Name: "java",
							Properties: Properties{
								"url": "${default-url}",
							},
						},
					},
				},
			}, {
				name: "GetModuleByName: Name don't exist ",
				fields: fields{
					Modules: []*Modules{
						{
							Name: "someproj-db",
							Type: "hdb",
							Path: "db",
							Requires: []Requires{
								{
									Name: "someproj-hdi-container",
								},
							},
						}}},
				moduleName: "foo",
				want:       nil,
				wantErr:    true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mta := &MTA{
					SchemaVersion: tt.fields.SchemaVersion,
					Id:            tt.fields.Id,
					Version:       tt.fields.Version,
					Modules:       tt.fields.Modules,
					Resources:     tt.fields.Resources,
					Parameters:    tt.fields.Parameters,
				}
				got, err := mta.GetModuleByName(tt.moduleName)
				if (err != nil) != tt.wantErr {
					t.Errorf("MTA.GetModuleByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MTA.GetModuleByName() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}
