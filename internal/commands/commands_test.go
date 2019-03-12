package commands

import (
	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"

	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Commands tests", func() {
	It("Mesh", func() {
		var moduleTypesCfg = []byte(`
version: 1
module-types:
  - name: html5
    info: "build UI application"
    commands:
    - command: npm install
    - command: grunt
    - command: npm prune --production
  - name: java
    info: "build java application"
    commands:
    - command: mvn clean install
  - name: nodejs
    info: "build nodejs application"
    commands:
    - command: npm install
  - name: golang
    info: "build golang application"
    commands:
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
		commands := ModuleTypes{}
		customCommands := Builders{}
		Ω(yaml.Unmarshal(moduleTypesCfg, &commands)).Should(Succeed())
		Ω(mesh(&modules, &commands, customCommands)).Should(Equal(expected))
		modules = mta.Module{
			Name: "uiapp1",
			Type: "html5",
			Path: "./",
		}
		_, err := mesh(&modules, &commands, customCommands)
		Ω(err).Should(Succeed())
		modules = mta.Module{
			Name: "uiapp1",
			Type: "html5",
			Path: "./",
			BuildParams: map[string]interface{}{
				"builder": "html5x",
			},
		}
		_, err = mesh(&modules, &commands, customCommands)
		Ω(err).Should(HaveOccurred())
	})

	It("Mesh - with builder in module types config", func() {
		var moduleTypesCfg = []byte(`
version: 1
module-types:
  - name: html5
    info: "build UI application"
    builder: npm
`)
		var buildersCfg = []byte(`
version: 1
builders:
  - name: npm
    info: "build UI application"
    commands:
    - command: npm install
    - command: npm prune --production
`)
		var modules = mta.Module{
			Name: "uiapp",
			Type: "html5",
			Path: "./",
		}
		var expected = CommandList{
			Info:    "build UI application",
			Command: []string{"npm install", "npm prune --production"},
		}
		commands := ModuleTypes{}
		customCommands := Builders{}
		Ω(yaml.Unmarshal(moduleTypesCfg, &commands)).Should(Succeed())
		Ω(yaml.Unmarshal(buildersCfg, &customCommands)).Should(Succeed())
		Ω(mesh(&modules, &commands, customCommands)).Should(Equal(expected))
	})

	It("Mesh - fails on usage both builder and commands in one module type", func() {
		var moduleTypesCfg = []byte(`
version: 1
module-types:
  - name: html5
    info: "build UI application"
    builder: npm
    commands:
    - command: npm install
    - command: grunt
    - command: npm prune --production
`)
		var modules = mta.Module{
			Name: "uiapp",
			Type: "html5",
			Path: "./",
		}
		commands := ModuleTypes{}
		customCommands := Builders{}
		Ω(yaml.Unmarshal(moduleTypesCfg, &commands)).Should(Succeed())
		_, err := mesh(&modules, &commands, customCommands)
		Ω(err).Should(HaveOccurred())
	})

	It("CommandProvider", func() {
		expected := CommandList{
			Info:    "installing module dependencies & execute grunt & remove dev dependencies",
			Command: []string{"npm install", "grunt", "npm prune --production"},
		}
		Ω(CommandProvider(mta.Module{Type: "html5"})).Should(Equal(expected))
	})

	var _ = Describe("CommandProvider - Invalid module types cfg", func() {
		var config []byte

		BeforeEach(func() {
			config = make([]byte, len(ModuleTypeConfig))
			copy(config, ModuleTypeConfig)
			// Simplified commands configuration (performance purposes). removed "npm prune --production"
			ModuleTypeConfig = []byte(`
module-types:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands: [xxx]
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  commands:
`)
		})

		AfterEach(func() {
			ModuleTypeConfig = make([]byte, len(config))
			copy(ModuleTypeConfig, config)
		})

		It("test", func() {
			_, err := CommandProvider(mta.Module{Type: "html5"})
			Ω(err).Should(HaveOccurred())
		})
	})

	var _ = Describe("CommandProvider - Invalid builders cfg", func() {
		var moduleTypesConfig []byte
		var buildersConfig []byte

		BeforeEach(func() {
			moduleTypesConfig = make([]byte, len(ModuleTypeConfig))
			copy(moduleTypesConfig, ModuleTypeConfig)
			// Simplified commands configuration (performance purposes). removed "npm prune --production"
			ModuleTypeConfig = []byte(`
module-types:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands:
`)

			buildersConfig = make([]byte, len(BuilderTypeConfig))
			copy(buildersConfig, BuilderTypeConfig)
			BuilderTypeConfig = []byte(`
builders:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands: [xxx]	
`)
		})

		AfterEach(func() {
			ModuleTypeConfig = make([]byte, len(moduleTypesConfig))
			copy(ModuleTypeConfig, moduleTypesConfig)
			BuilderTypeConfig = make([]byte, len(buildersConfig))
			copy(BuilderTypeConfig, buildersConfig)
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
		It("Fetcher builder specified in build params", func() {
			var mtaCF = []byte(`
_schema-version: "2.0.0"
ID: mta_proj
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
    build-parameters:
      builder: fetcher
      fetcher-opts:
         repo-type: maven
         repo-coordinates: com.sap.xs.java:xs-audit-log-api:1.2.3

`)
			m := mta.MTA{}
			// parse mta yaml
			Ω(yaml.Unmarshal(mtaCF, &m)).Should(Succeed())
			module, commands, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(BeNil())
			Ω(module.Path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"mvn dependency:copy -Dartifact=com.sap.xs.java:xs-audit-log-api:1.2.3 -DoutputDirectory=./"}))
		})
	})

	It("Invalid case - wrong fetcher builder type", func() {
		wd, _ := os.Getwd()
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtaWithFetcher.yaml"}
		_, _, err := GetModuleAndCommands(&ep, "j1")
		Ω(err).Should(HaveOccurred())
	})

	It("Invalid case - wrong mta", func() {
		wd, _ := os.Getwd()
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtaUnknown.yaml"}
		_, _, err := GetModuleAndCommands(&ep, "node-js")
		Ω(err).Should(HaveOccurred())

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
			conf := ModuleTypeConfig
			ModuleTypeConfig = []byte("wrong config")
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			_, _, err := GetModuleAndCommands(&ep, "node-js")
			ModuleTypeConfig = conf
			Ω(err).Should(HaveOccurred())
		})
	})
})
