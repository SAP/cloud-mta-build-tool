package builders

import (
	"cloud-mta-build-tool/mta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Commands tests", func() {
	It("Mesh", func() {
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
		var modules = mta.Modules{
			Name: "uiapp",
			Type: "html5",
			Path: "./",
		}
		var expected = commandList{
			Info:    "build UI application",
			Command: []string{"npm install", "grunt", "npm prune --production"},
		}
		commands := Builders{}
		Ω(yaml.Unmarshal(buildersCfg, &commands)).Should(Succeed())
		Ω(mesh(modules, commands)).Should(Equal(expected))
	})

	It("CommandProvider", func() {
		expected := commandList{
			Info:    "installing module dependencies & execute grunt & remove dev dependencies",
			Command: []string{"npm install", "grunt", "npm prune --production"},
		}
		Ω(CommandProvider(mta.Modules{Type: "html5"})).Should(Equal(expected))
	})
})
