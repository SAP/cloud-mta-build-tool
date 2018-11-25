package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Execute", func() {

	var _ = Describe("Execute call", func() {

		var _ = DescribeTable("Valid input", func(args [][]string) {
			Ω(Execute(args)).Should(Succeed())
		},
			Entry("EchoTesting", [][]string{{"", "echo", "-n", `{"Name": "Bob", "Age": 32}`}}),
			Entry("Dummy Go Testing", [][]string{{"", "go", "test", "exec_dummy_test.go"}}))

		var _ = DescribeTable("Invalid input", func(args [][]string) {
			Ω(Execute(args)).Should(HaveOccurred())
		},
			Entry("Valid command fails on input", [][]string{{"", "go", "test", "exec_unknown_test.go"}}),
			Entry("Invalid command", [][]string{{"", "dateXXX"}}),
		)
	})

	It("Indicator", func() {
		// var wg sync.WaitGroup
		// wg.Add(1)
		shutdownCh := make(chan struct{})
		start := time.Now()
		go indicator(shutdownCh)
		time.Sleep(3 * time.Second)
		// close(shutdownCh)
		sec := time.Since(start).Seconds()
		switch int(sec) {
		case 0:
			// Output:
		case 1:
			// Output: .
		case 2:
			// Output: ..
		case 3:
			// Output: ...
		default:
		}

		shutdownCh <- struct{}{}
		// wg.Wait()
	})

	var _ = Describe("copyModuleArchive", func() {

		AfterEach(func() {
			os.Remove(getTestPath("mta", "node-js", "data.zip"))
			os.RemoveAll(getTestPath("mta", "mta"))
			mta.GetWorkingDirectory = mta.OsGetWd
		})

		It("Sanity", func() {
			lp := mta.Loc{SourcePath: getTestPath("mta")}
			dir.Archive(getTestPath("mta", "node-js"), getTestPath("mta", "node-js", "data.zip"))

			Ω(copyModuleArchive(&lp, "node-js", "node-js")).Should(Succeed())
			Ω(getTestPath("mta", "mta", "node-js", "data.zip")).Should(BeAnExistingFile())
		})

		var _ = DescribeTable("Invalid cases", func(modulePath string, mockWd bool, failOnCall int) {
			var lp mta.Loc
			var countCalls = 0
			if mockWd {
				lp = mta.Loc{}
				mta.GetWorkingDirectory = func() (string, error) {
					countCalls++
					if countCalls >= failOnCall {
						return "", errors.New("error")
					}
					return os.Getwd()
				}
			} else {
				lp = mta.Loc{SourcePath: getTestPath("mta")}
			}
			Ω(copyModuleArchive(&lp, modulePath, "node-js")).Should(HaveOccurred())
		},
			Entry("Invalid module name", "node-js1", false, -1),
			Entry("Get wd fails", "node-js", true, 1),
			Entry("Get wd fails", "node-js", true, 2))
	})

	var _ = Describe("Validation", func() {
		var _ = DescribeTable("getValidationMode", func(flag string, expectedValidateSchema, expectedValidateProject, expectedSuccess bool) {
			res1, res2, err := GetValidationMode(flag)
			Ω(res1).Should(Equal(expectedValidateSchema))
			Ω(res2).Should(Equal(expectedValidateProject))
			Ω(err == nil).Should(Equal(expectedSuccess))
		},
			Entry("all", "", true, true, true),
			Entry("schema", "schema", true, false, true),
			Entry("project", "project", false, true, true),
			Entry("invalid", "value", false, false, false),
		)

		var _ = DescribeTable("validateMtaYaml", func(projectRelPath string, validateSchema, validateProject, expectedSuccess bool) {
			ep := mta.Loc{SourcePath: getTestPath(projectRelPath)}
			err := ValidateMtaYaml(&ep, validateSchema, validateProject)
			Ω(err == nil).Should(Equal(expectedSuccess))
		},
			Entry("invalid path to yaml - all", "ui5app1", true, true, false),
			Entry("invalid path to yaml - schema", "ui5app1", true, false, false),
			Entry("invalid path to yaml - project", "ui5app1", false, true, false),
			Entry("invalid path to yaml - nothing to validate", "ui5app1", false, false, true),
			Entry("valid schema", "mtahtml5", true, false, true),
			Entry("invalid project - no ui5app2 path", "mtahtml5", false, true, false),
		)
	})
})

var _ = Describe("Build", func() {

	var _ = Describe("build Module", func() {

		var config []byte

		BeforeEach(func() {
			config = make([]byte, len(builders.CommandsConfig))
			copy(config, builders.CommandsConfig)
			// Simplified commands configuration (performance purposes). removed "npm prune --production"
			builders.CommandsConfig = []byte(`
builders:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  type:
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  type:
`)
		})

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
			builders.CommandsConfig = make([]byte, len(config))
			copy(builders.CommandsConfig, config)
		})

		It("Sanity", func() {
			ep := mta.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			Ω(BuildModule(&ep, "node-js")).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})

		var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
			ep := mta.Loc{SourcePath: getTestPath(projectName), TargetPath: getTestPath("result"), MtaFilename: mtaFilename}
			Ω(BuildModule(&ep, moduleName)).Should(HaveOccurred())
			Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
		},
			Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
			Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
			Entry("Invalid module name", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
		)

	})

	var _ = Describe("Generate Commands", func() {

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		readFileContent := func(filename string) string {
			content, _ := ioutil.ReadFile(filename)
			contentString := string(content[:])
			contentString = strings.Replace(contentString, "\n", "", -1)
			contentString = strings.Replace(contentString, "\r", "", -1)
			return contentString
		}

		It("Generate Meta", func() {
			ep := mta.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			GenerateMeta(&ep)
			mtadPath, _ := ep.GetMtadPath()
			Ω(readFileContent(mtadPath)).Should(Equal(readFileContent(getTestPath("golden", "mtad.yaml"))))
		})

		It("Generate Mtar", func() {
			ep := mta.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			err := GenerateMeta(&ep)
			if err != nil {
				fmt.Println(err)
			}
			err = GenerateMtar(&ep)
			if err != nil {
				fmt.Println(err)
			}
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})

	var _ = Describe("Command converter", func() {

		It("Sanity", func() {
			cmdInput := []string{"npm install", "grunt", "npm prune --production"}
			cmdExpected := [][]string{
				{"path", "npm", "install"},
				{"path", "grunt"},
				{"path", "npm", "prune", "--production"}}
			Ω(cmdConverter("path", cmdInput)).Should(Equal(cmdExpected))
		})
	})

	var _ = Describe("Process Dependencies", func() {
		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		It("Sanity", func() {
			Ω(processDependencies(&mta.Loc{SourcePath: getTestPath("mtahtml5"), MtaFilename: "mtaWithBuildParams.yaml"}, "ui5app")).Should(Succeed())
		})
		It("Invalid mta", func() {
			Ω(processDependencies(&mta.Loc{SourcePath: getTestPath("mtahtml5"), MtaFilename: "mta1.yaml"}, "ui5app")).Should(HaveOccurred())
		})
		It("Invalid module name", func() {
			Ω(processDependencies(&mta.Loc{SourcePath: getTestPath("mtahtml5")}, "xxx")).Should(HaveOccurred())
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
			err := yaml.Unmarshal(mtaCF, &m)
			if err != nil {
				fmt.Println(err)
			}
			module, commands, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(BeNil())
			Ω(module.Path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"npm install", "grunt", "npm prune --production"}))
		})
	})
})

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}
