package mta

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"cloud-mta-build-tool/cmd/mta/models"
	"github.com/stretchr/testify/assert"
)

type validator = func(t *testing.T, actual, expected models.Modules)

type testInfo struct {
	name     string
	expected models.Modules
}

func doTest(t *testing.T, expected []testInfo, validators []validator, filename string) {
	mtaFile, _ := ioutil.ReadFile(filename)

	actual, _ := Parse(mtaFile)
	for i, tt := range expected {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, actual)
			require.Len(t, actual.Modules, len(expected))
			validators[i](t, *actual.Modules[i], tt.expected)
		})
	}
	mtaContent, err := Marshal(actual)
	assert.Nil(t, err)
	newActual, newErr := Parse(mtaContent)
	assert.Nil(t, newErr)
	assert.Equal(t, actual, newActual)
}

// Table driven test
// Unit test for parsing mta files to working object
func Test_ModulesParsing(t *testing.T) {
	tests := []testInfo{
		{
			name: "Parse service(srv) Module section",
			expected: models.Modules{
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
				Properties: models.Properties{
					"APPC_LOG_LEVEL":              "info",
					"VSCODE_JAVA_DEBUG_LOG_LEVEL": "ALL",
				},
			},
		},

		// ------------------------Second module test------------------------------
		{
			name: "Parse UI(HTML5) Module section",
			expected: models.Modules{
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
					Builder: "grunt",
				},

				Parameters: models.Parameters{
					"disk-quota": "256M",
					"memory":     "256M",
				},
			},
		},
	}

	validators := [2]validator{
		func(t *testing.T, actual, expected models.Modules) {
			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expected.Type, actual.Type)
			assert.Equal(t, expected.Path, actual.Path)
			assert.Equal(t, expected.Parameters, actual.Parameters)
			assert.Equal(t, expected.Properties, actual.Properties)
			assert.Equal(t, expected.Requires, actual.Requires)
			assert.Equal(t, expected.Provides, actual.Provides)
		},
		func(t *testing.T, actual, expected models.Modules) {
			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expected.Type, actual.Type)
			assert.Equal(t, expected.Path, actual.Path)
			assert.Equal(t, expected.Requires, actual.Requires)
			assert.Equal(t, expected.Parameters, actual.Parameters)
			assert.Equal(t, expected.BuildParams, actual.BuildParams)
		},
	}

	doTest(t, tests, validators[:], "./testdata/mta.yaml")

}

func Test_BrokenMta(t *testing.T){
	mtaContent, _ := ioutil.ReadFile("./testdata/mtaWithBrokenProperties.yaml")

	mta, err := Parse(mtaContent)
	require.NotNil(t, err)
	require.NotNil(t, mta)
}

func Test_FullMta(t *testing.T) {
	schemaVersion := "2.0.0"

	expected := models.MTA{
		SchemaVersion: &schemaVersion,
		Id:            "cloud.samples.someproj",
		Version:       "1.0.0",
		Parameters: models.Parameters{
			"deploy_mode": "html5-repo",
		},
		Modules: []*models.Modules{
			{
				Name: "someproj-db",
				Type: "hdb",
				Path: "db",
				Requires: []models.Requires{
					{
						Name: "someproj-hdi-container",
					},
					{
						Name: "someproj-logging",
					},
				},
				Parameters: models.Parameters{
					"disk-quota": "256M",
					"memory":     "256M",
				},
			},
			{
				Name: "someproj-java",
				Type: "java",
				Path: "srv",
				Parameters: models.Parameters{
					"memory":     "512M",
					"disk-quota": "256M",
				},
				Provides: []models.Provides{
					{
						Name: "java",
						Properties: models.Properties{
							"url": "${default-url}",
						},
					},
				},
				Requires: []models.Requires{
					{
						Name: "someproj-hdi-container",
						Properties: models.Properties{
							"JBP_CONFIG_RESOURCE_CONFIGURATION":
							"[tomcat/webapps/ROOT/META-INF/context.xml: " +
								"{\"service_name_for_DefaultDB\" : \"~{hdi-container-name}\"}]",
						},
					},
					{
						Name: "someproj-logging",
					},
				},
				BuildParams: models.BuildParameters{
					Requires: []models.BuildRequires{
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
				Parameters: models.Parameters{
					"memory":     "256M",
					"disk-quota": "256M",
				},
				Requires: []models.Requires{
					{
						Name:  "java",
						Group: "destinations",
						Properties: models.Properties{
							"name": "someproj-backend",
							"url":  "~{url}",
						},
					},
					{
						Name: "someproj-logging",
					},
				},
				BuildParams: models.BuildParameters{
					Builder: "grunt",
					Requires: []models.BuildRequires{
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
				Parameters: models.Parameters{
					"memory":     "256M",
					"disk-quota": "256M",
				},
				Requires: []models.Requires{
					{
						Name: "someproj-apprepo-dt",
					},
				},
				BuildParams: models.BuildParameters{
					Builder: "grunt",
					Type:    "com.sap.html5.application-content",
					Requires: []models.BuildRequires{
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
				Parameters: models.Parameters{
					"memory":     "256M",
					"disk-quota": "256M",
				},
				Requires: []models.Requires{
					{
						Name:  "java",
						Group: "destinations",
						Properties: models.Properties{
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
		Resources: []*models.Resources{
			{
				Name: "someproj-hdi-container",
				Properties: models.Properties{
					"hdi-container-name": "${service-name}",
				},
				Type: "com.sap.xs.hdi-container",
			},
			{
				Name: "someproj-apprepo-rt",
				Type: "org.cloudfoundry.managed-service",
				Parameters: models.Parameters{
					"service":      "html5-apps-repo",
					"service-plan": "app-runtime",
				},
			},
			{
				Name: "someproj-apprepo-dt",
				Type: "org.cloudfoundry.managed-service",
				Parameters: models.Parameters{
					"service":      "html5-apps-repo",
					"service-plan": "app-host",
				},
			},
			{
				Name: "someproj-logging",
				Type: "org.cloudfoundry.managed-service",
				Parameters: models.Parameters{
					"service":      "application-logs",
					"service-plan": "lite",
				},
			},
		},
	}

	mtaContent, _ := ioutil.ReadFile("./testdata/mta2.yaml")

	actual, err := Parse(mtaContent)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

}
