package commands

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
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
		Ω(mesh(&modules, &commands, &customCommands)).Should(Equal(expected))
		modules = mta.Module{
			Name: "uiapp1",
			Type: "html5",
			Path: "./",
		}
		_, _, err := mesh(&modules, &commands, &customCommands)
		Ω(err).Should(Succeed())
		modules = mta.Module{
			Name: "uiapp1",
			Type: "html5",
			Path: "./",
			BuildParams: map[string]interface{}{
				"builder": "html5x",
			},
		}
		_, _, err = mesh(&modules, &commands, &customCommands)
		Ω(err).Should(HaveOccurred())
	})

	It("Mesh - fails on wrong type of property commands", func() {
		module := mta.Module{
			Name: "uiapp1",
			Type: "html5",
			Path: "./",
			BuildParams: map[string]interface{}{
				"builder":  "custom",
				"commands": "cmd",
			},
		}
		commands := ModuleTypes{}
		customCommands := Builders{}
		_, _, err := mesh(&module, &commands, &customCommands)
		Ω(err).Should(HaveOccurred())
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
		_, _, err := mesh(&modules, &commands, &customCommands)
		Ω(err).Should(HaveOccurred())
	})

	It("Mesh - custom builder", func() {
		var modules = mta.Module{
			Name: "uiapp",
			Type: "html5",
			Path: "./",
			BuildParams: map[string]interface{}{
				builderParam:  "custom",
				commandsParam: []string{"command1"},
			},
		}
		commands := ModuleTypes{}
		customCommands := Builders{}
		cmds, _, err := mesh(&modules, &commands, &customCommands)
		Ω(err).Should(Succeed())
		Ω(len(cmds.Command)).Should(Equal(1))
		Ω(cmds.Command[0]).Should(Equal("command1"))
	})

	It("CommandProvider", func() {
		expected := CommandList{
			Info:    "installing module dependencies & remove dev dependencies",
			Command: []string{"npm install --production"},
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
			_, _, err := CommandProvider(mta.Module{Type: "html5"})
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
			_, _, err := CommandProvider(mta.Module{Type: "html5"})
			Ω(err).Should(HaveOccurred())
		})
	})

	var _ = Describe("Command converter", func() {

		It("Sanity", func() {
			cmdInput := []string{"npm install {{config}}", "grunt", "npm prune --production"}
			cmdExpected := [][]string{
				{"path", "npm", "install", "{{config}}"},
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
			module, commands, _, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(Succeed())
			Ω(module.Path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"npm install --production"}))
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
			module, commands, _, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(BeNil())
			Ω(module.Path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"npm install --production"}))
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
			module, commands, _, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(BeNil())
			Ω(module.Path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"mvn -B dependency:copy -Dartifact=com.sap.xs.java:xs-audit-log-api:1.2.3 -DoutputDirectory=./target"}))
		})

		It("Invalid case - wrong mta", func() {
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtaUnknown.yaml"}
			_, _, _, err := GetModuleAndCommands(&ep, "node-js")
			Ω(err).Should(HaveOccurred())

		})

		var _ = Describe("GetModuleAndCommands", func() {
			It("Sanity", func() {
				wd, _ := os.Getwd()
				ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
				module, cmd, _, err := GetModuleAndCommands(&ep, "node-js")
				Ω(err).Should(Succeed())
				Ω(module.Name).Should(Equal("node-js"))
				Ω(len(cmd)).Should(Equal(1))
				Ω(cmd[0]).Should(Equal("npm install --production"))

			})
			It("Invalid case - wrong module name", func() {
				wd, _ := os.Getwd()
				ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
				_, _, _, err := GetModuleAndCommands(&ep, "node-js1")
				Ω(err).Should(HaveOccurred())

			})
			It("Invalid case - wrong mta", func() {
				wd, _ := os.Getwd()
				ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtaUnknown.yaml"}
				_, _, _, err := GetModuleAndCommands(&ep, "node-js")
				Ω(err).Should(HaveOccurred())

			})
			It("Invalid case - wrong type", func() {
				wd, _ := os.Getwd()
				ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtaUnknownBuilder.yaml"}
				_, cmd, _, _ := GetModuleAndCommands(&ep, "node-js")
				Ω(len(cmd)).Should(Equal(0))

			})
			It("Invalid case - broken commands config", func() {
				conf := ModuleTypeConfig
				ModuleTypeConfig = []byte("wrong config")
				wd, _ := os.Getwd()
				ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
				_, _, _, err := GetModuleAndCommands(&ep, "node-js")
				ModuleTypeConfig = conf
				Ω(err).Should(HaveOccurred())
			})
		})
	})

	var _ = Describe("GetBuilder", func() {
		It("Builder defined by type", func() {
			m := mta.Module{
				Name: "x",
				Type: "node-js",
			}
			Ω(GetBuilder(&m)).Should(Equal("node-js"))
		})
		It("Builder defined by build params", func() {
			m := mta.Module{
				Name: "x",
				Type: "node-js",
				BuildParams: map[string]interface{}{
					builderParam: "npm",
				},
			}
			builder, custom, cmds, _, err := GetBuilder(&m)
			Ω(builder).Should(Equal("npm"))
			Ω(custom).Should(Equal(true))
			Ω(len(cmds)).Should(Equal(0))
			Ω(err).Should(Succeed())
		})
		It("Custom builder with no commands", func() {
			m := mta.Module{
				Name: "x",
				Type: "node-js",
				BuildParams: map[string]interface{}{
					builderParam: customBuilder,
				},
			}
			builder, custom, _, _, err := GetBuilder(&m)
			Ω(builder).Should(Equal(customBuilder))
			Ω(custom).Should(Equal(true))
			Ω(err).Should(Succeed())
		})
		It("Custom builder with wrong commands definition", func() {
			m := mta.Module{
				Name: "x",
				Type: "node-js",
				BuildParams: map[string]interface{}{
					builderParam:  customBuilder,
					commandsParam: "command1",
				},
			}
			builder, custom, _, _, err := GetBuilder(&m)
			Ω(builder).Should(Equal(customBuilder))
			Ω(custom).Should(Equal(true))
			Ω(err.Error()).Should(Equal(`the "commands" property is defined incorrectly; the property must contain a sequence of strings`))
		})
	})
})
