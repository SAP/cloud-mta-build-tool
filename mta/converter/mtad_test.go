package converter

import (
	"log"
	"reflect"
	"testing"

	"cloud-mta-build-tool/cmd/platform"
	"cloud-mta-build-tool/mta/models"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestConvertTypes(t *testing.T) {

	// ------------platform Config content ---------
	var platformYaml = []byte(`
platform:
  - name: cf
    modules:
    - native-type: html5
      platform-type: "javascript.nodejs"
    - native-type: nodejs
      platform-type: "javascript.nodejs"
    - native-type: java
      platform-type: "java.tomcat"
    - native-type: hdb
      platform-type: "com.sap.xs.hdi"

  - name: neo
    modules:
    - native-type: html5
      platform-type: "com.sap.hcp.html5"
    - native-type: java
      platform-type: "java.tomcat"
`)

	platformType := platform.Platforms{}
	// parse mta yaml
	err := yaml.Unmarshal(platformYaml, &platformType)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)
	}

	//---------------- MTA single module content-------------------------
	var mtaSingleModule = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
`)

	mta := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(mtaSingleModule, &mta)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	//expected one module
	var expectedMta1Modules = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: javascript.nodejs
    path: app
`)

	expected := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(expectedMta1Modules, &expected)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)

	}

	//----------------Multi Neo--------------

	// MTA content
	var mtaNeo = []byte(`
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

	//Parse the mta content
	actualMtaMultiNeo := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(mtaNeo, &actualMtaMultiNeo)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	//expected for multi Neo
	var expectedMtaMultiModules = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: com.sap.hcp.html5
    path: app

  - name: htmlapp2
    type: com.sap.hcp.html5
    path: app

  - name: java
    type: java.tomcat
    path: app
`)

	//Parse the expected content
	expectedMultiModulesNeo := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(expectedMtaMultiModules, &expectedMultiModulesNeo)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)

	}

	//----------------Multi CF----------------------------------
	// MTA content
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

	actulMtaCFMulti := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(mtaCF, &actulMtaCFMulti)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	//expected for multi modules
	var expectedMultiModCF = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: com.sap.hcp.html5
    path: app

  - name: htmlapp2
    type: com.sap.hcp.html5
    path: app

  - name: java
    type: java.tomcat
    path: app
`)

	expectedMultiModulesCF := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(expectedMultiModCF, &expectedMultiModulesCF)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	tests := []struct {
		name          string
		mta           models.MTA
		platforms     platform.Platforms
		platform      string
		expected      string
		expectedMulti models.MTA
	}{
		{

			name:      "Module with one platform config",
			mta:       mta,
			platforms: platformType,
			platform:  "cf",
			expected:  expected.Modules[0].Type,
		},
		{
			name:          "Multi modules multi platforms config Neo",
			mta:           actualMtaMultiNeo,
			platforms:     platformType,
			platform:      "neo",
			expectedMulti: expectedMultiModulesNeo,
		},
		{
			name:          "Multi modules multi platforms config CF",
			mta:           actulMtaCFMulti,
			platforms:     platformType,
			platform:      "cf",
			expectedMulti: expectedMultiModulesCF,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case 0:
				//One module convert
				ConvertTypes(tt.mta, tt.platforms, tt.platform)
				if !assert.Equal(t, mta.Modules[0].Type, tt.expected) {
					t.Error("Test was failed")
				}
			case 1:
				//Multi module convert neo
				ConvertTypes(tt.mta, tt.platforms, tt.platform)
				if !assert.Equal(t, actualMtaMultiNeo.Modules, tt.expectedMulti.Modules) {
					t.Error("Test was failed")
				}
			case 2:
				//Multi module convert cloud foundry
				ConvertTypes(tt.mta, tt.platforms, tt.platform)
				if !assert.Equal(t, actualMtaMultiNeo.Modules, tt.expectedMulti.Modules) {
					t.Error("Test was failed")
				}
			}
		})
	}
}

func TestPlatformConfig(t *testing.T) {

	// ------------Multi platform ---------
	var platformsCfgMulti = []byte(`
platform:
  - name: cf
    modules:
    - native-type: html5
      platform-type: "javascript.nodejs"
    - native-type: nodejs
      platform-type: "javascript.nodejs"
    - native-type: java
      platform-type: "java.tomcat"
    - native-type: hdb
      platform-type: "com.sap.xs.hdi"

  - name: neo
    modules:
    - native-type: html5
      platform-type: "com.sap.hcp.html5"
    - native-type: java
      platform-type: "java.tomcat"
`)

	pl := platform.Platforms{}
	// parse mta yaml
	err := yaml.Unmarshal(platformsCfgMulti, &pl)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)
	}

	// ------------One platform ---------
	var platformsCfgSingle = []byte(`
platform:
  - name: cf
    modules:
    - native-type: html5
      platform-type: "javascript.nodejs"
    - native-type: nodejs
      platform-type: "javascript.nodejs"
    - native-type: java
      platform-type: "java.tomcat"
    - native-type: hdb
      platform-type: "com.sap.xs.hdi"
`)

	ps := platform.Platforms{}
	// parse mta yaml
	err = yaml.Unmarshal(platformsCfgSingle, &ps)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)
	}

	tests := []struct {
		name      string
		platform  string
		platforms platform.Platforms
		expected  platform.Platforms
	}{

		{
			name:      "Platform test",
			platform:  "cf",
			platforms: pl,
			expected:  ps,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PlatformConfig(tt.platforms, tt.platform)
			if !reflect.DeepEqual(got, tt.expected.Platforms[0]) {
				t.Errorf("PlatformConfig() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
