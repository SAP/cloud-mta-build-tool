package builders

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/mta"
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
		var modules = mta.Module{
			Name: "uiapp",
			Type: "html5",
			Path: "./",
		}
		var expected = CommandList{
			Info:    "build UI application",
			Command: []string{"npm install", "grunt", "npm prune --production"},
		}
		commands := Builders{}
		Ω(yaml.Unmarshal(buildersCfg, &commands)).Should(Succeed())
		Ω(mesh(modules, commands)).Should(Equal(expected))
	})

	It("CommandProvider", func() {
		expected := CommandList{
			Info:    "installing module dependencies & execute grunt & remove dev dependencies",
			Command: []string{"npm install", "grunt", "npm prune --production"},
		}
		Ω(CommandProvider(mta.Module{Type: "html5"})).Should(Equal(expected))
	})

	var _ = Describe("CommandProvider - Invalid cfg", func() {
		var config []byte

		BeforeEach(func() {
			config = make([]byte, len(CommandsConfig))
			copy(config, CommandsConfig)
			// Simplified commands configuration (performance purposes). removed "npm prune --production"
			CommandsConfig = []byte(`
builders:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  type: [xxx]
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  type:
`)
		})

		AfterEach(func() {
			CommandsConfig = make([]byte, len(config))
			copy(CommandsConfig, config)
		})

		It("test", func() {
			_, err := CommandProvider(mta.Module{Type: "html5"})
			Ω(err).Should(HaveOccurred())
		})
	})
})
