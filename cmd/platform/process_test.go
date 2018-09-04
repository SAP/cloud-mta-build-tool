package platform

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestParse(t *testing.T) {

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
