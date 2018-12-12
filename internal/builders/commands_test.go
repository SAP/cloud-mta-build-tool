package builders

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fs"

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

	var _ = Describe("Command converter", func() {

		It("Sanity", func() {
			cmdInput := []string{"npm install", "grunt", "npm prune --production"}
			cmdExpected := [][]string{
				{"path", "npm", "install"},
				{"path", "grunt"},
				{"path", "npm", "prune", "--production"}}
			Ω(CmdConverter("path", cmdInput)).Should(Equal(cmdExpected))
		})
	})

	var _ = Describe("moduleCmd", func() {
		It("Sanity", func() {
			var mtaCF = []byte(`
_schema-version: "2.0.0"
ID: mta_proj
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

			m := mta.MTA{}
			// parse mta yaml
			Ω(yaml.Unmarshal(mtaCF, &m)).Should(Succeed())
			module, commands, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(Succeed())
			Ω(module.Path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"npm install", "grunt", "npm prune --production"}))
		})

		It("Builder specified in build params", func() {
			var mtaCF = []byte(`
_schema-version: "2.0.0"
ID: mta_proj
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
    build-parameters:
      builder: npm
`)

			m := mta.MTA{}
			// parse mta yaml
			Ω(yaml.Unmarshal(mtaCF, &m)).Should(Succeed())
			module, commands, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(BeNil())
			Ω(module.Path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"npm install", "npm prune --production"}))
		})
	})

	var _ = Describe("GetModuleAndCommands", func() {
		It("Sanity", func() {
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			module, cmd, err := GetModuleAndCommands(&ep, "node-js")
			Ω(err).Should(Succeed())
			Ω(module.Name).Should(Equal("node-js"))
			Ω(len(cmd)).Should(Equal(1))
			Ω(cmd[0]).Should(Equal("npm prune --production"))

		})
		It("Invalid case - wrong module name", func() {
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			_, _, err := GetModuleAndCommands(&ep, "node-js1")
			Ω(err).Should(HaveOccurred())

		})
		It("Invalid case - wrong mta", func() {
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtaUnknown.yaml"}
			_, _, err := GetModuleAndCommands(&ep, "node-js")
			Ω(err).Should(HaveOccurred())

		})
		It("Invalid case - wrong type", func() {
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtaUnknownBuilder.yaml"}
			_, cmd, _ := GetModuleAndCommands(&ep, "node-js")
			Ω(len(cmd)).Should(Equal(0))

		})
		It("Invalid case - broken commands config", func() {
			conf := CommandsConfig
			CommandsConfig = []byte("wrong config")
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			_, _, err := GetModuleAndCommands(&ep, "node-js")
			CommandsConfig = conf
			Ω(err).Should(HaveOccurred())
		})
	})
})
