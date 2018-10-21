package platform

import (
	"log"
	"reflect"
	"testing"

	"cloud-mta-build-tool/mta"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestParse(t *testing.T) {
	t.Parallel()
	// ------------Multi platform ---------
	var platforms = []byte(`
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
     platform-type: "com.sap.hcp.html"
   - native-type: java
     platform-type: "java.tomcat"
`)

	ps := Platforms{}
	err := yaml.Unmarshal(platforms, &ps)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)
	}

	tests := []struct {
		name      string
		platforms []byte
		expected  Platforms
	}{
		{
			name:      "Platform test",
			platforms: platforms,
			expected:  ps,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.platforms)
			if !assert.Equal(t, got.Platforms, tt.expected.Platforms) {
				t.Errorf("Parse() = %v, `\n` want %v", got, tt.expected)
			}
		})
	}
}

func TestConvertType(t *testing.T) {

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

	platformType := Platforms{}
	// parse mta yaml
	err := yaml.Unmarshal(platformYaml, &platformType)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)
	}

	// ---------------- MTA single module content-------------------------
	var mtaSingleModule = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
`)

	m := mta.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(mtaSingleModule, &m)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	// expected one module
	var expectedMta1Modules = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: javascript.nodejs
    path: app
`)

	expected := mta.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(expectedMta1Modules, &expected)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)

	}

	// ----------------Multi Neo--------------

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

	// Parse the mta content
	actualMtaMultiNeo := mta.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(mtaNeo, &actualMtaMultiNeo)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	// expected for multi Neo
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

	// Parse the expected content
	expectedMultiModulesNeo := mta.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(expectedMtaMultiModules, &expectedMultiModulesNeo)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)

	}

	// ----------------Multi CF----------------------------------
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

	actulMtaCFMulti := mta.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(mtaCF, &actulMtaCFMulti)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	// expected for multi modules
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

	expectedMultiModulesCF := mta.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(expectedMultiModCF, &expectedMultiModulesCF)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	tests := []struct {
		name          string
		mta           mta.MTA
		platforms     Platforms
		platform      string
		expected      string
		expectedMulti mta.MTA
	}{
		{

			name:      "Module with one platform config",
			mta:       m,
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
				// One module convert
				ConvertTypes(tt.mta, tt.platforms, tt.platform)
				if !assert.Equal(t, m.Modules[0].Type, tt.expected) {
					t.Error("Test was failed")
				}
			case 1:
				// Multi module convert neo
				ConvertTypes(tt.mta, tt.platforms, tt.platform)
				if !assert.Equal(t, actualMtaMultiNeo.Modules, tt.expectedMulti.Modules) {
					t.Error("Test was failed")
				}
			case 2:
				// Multi module convert cloud foundry
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

	pl := Platforms{}
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

	ps := Platforms{}
	// parse mta yaml
	err = yaml.Unmarshal(platformsCfgSingle, &ps)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)
	}

	tests := []struct {
		name      string
		platform  string
		platforms Platforms
		expected  Platforms
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
			got := platformConfig(tt.platforms, tt.platform)
			if !reflect.DeepEqual(got, tt.expected.Platforms[0]) {
				t.Errorf("platformConfig() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
