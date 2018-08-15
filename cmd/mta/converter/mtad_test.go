package converter

import (
	"log"
	"testing"

	"cloud-mta-build-tool/cmd/mta/models"
	"cloud-mta-build-tool/cmd/platform"

	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
)

func TestConvertTypes(t *testing.T) {

	// MTA content
	var mtaYaml = []byte(`
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
	err := yaml.Unmarshal(mtaYaml, &mta)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}

	// Config content
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
	err = yaml.Unmarshal(platformYaml, &platformType)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)
	}

	//expected for 1 module
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
	err = yaml.Unmarshal(expectedMta1Modules , &expected)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)

	}


	//expected for 1 module
	var expectedMtaMultiModules = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: com.sap.hcp.html5
    path: app

  - name: htmlapp2
    type: java.tomcat
    path: app
`)


	expectedMultiModules := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(expectedMtaMultiModules , &expectedMultiModules)
	if err != nil {
		log.Fatalf("Error to parse platform yaml: %v", err)

	}


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
    type: java
    path: app
`)

	mtaNeoMulti  := models.MTA{}
	// parse mta yaml
	err = yaml.Unmarshal(mtaNeo, &mtaNeoMulti)
	if err != nil {
		log.Fatalf("Error to parse mta yaml: %v", err)
	}


	tests := []struct {
		name      string
		mta       models.MTA
		platforms platform.Platforms
		platform  string
		expected  string
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
			name:      "Multi modules multi platforms config",
			mta:       mtaNeoMulti,
			platforms: platformType,
			platform:  "neo",
			expectedMulti:  expectedMultiModules,
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
				//Multi module convert
				ConvertTypes(tt.mta, tt.platforms, tt.platform)
				if !assert.Equal(t, mtaNeoMulti.Modules, tt.expectedMulti.Modules) {
					t.Error("Test was failed")
				}
			}
		})
	}
}
