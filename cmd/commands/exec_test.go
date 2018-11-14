package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("Commands", func() {

	BeforeEach(func() {
		targetMtadFlag = getTestPath("result")
		targetMetaFlag = getTestPath("result")
		targetMtarFlag = getTestPath("result")
		targetPackFlag = getTestPath("result")
		targetBModuleFlag = getTestPath("result")
		targetCleanupFlag = getTestPath("result")
		logs.Logger = logs.NewLogger()
		err := os.Mkdir(targetMtadFlag, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterEach(func() {
		os.RemoveAll(targetMtadFlag)
	})

	var _ = Describe("Pack and cleanup commands", func() {
		It("Target file in opened status", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			// Target path has to be dir, but is currently created and opened as file
			pPackModuleFlag = "ui5app"
			sourcePackFlag = getTestPath("mtahtml5")
			ep := dir.MtaLocationParameters{SourcePath: sourcePackFlag, TargetPath: targetPackFlag}
			targetTmpDir, _ := ep.GetTargetTmpDir()
			err := os.MkdirAll(targetTmpDir, os.ModePerm)
			if err != nil {
				logs.Logger.Error(err)
			}
			f, _ := os.Create(filepath.Join(targetTmpDir, "ui5app"))
			packCmd.Run(nil, []string{})
			fmt.Println(str.String())
			Ω(str.String()).Should(ContainSubstring("Pack of module ui5app failed on making directory"))

			f.Close()
			// cleanup command used for test temp file removal
			sourceCleanupFlag = sourcePackFlag
			cleanupCmd.Run(nil, []string{})
			Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
		})

	})

	var _ = Describe("Generate commands call", func() {

		var ep dir.MtaLocationParameters

		It("Generate Meta", func() {
			sourceMetaFlag = getTestPath("mtahtml5")
			ep = dir.MtaLocationParameters{SourcePath: sourceMetaFlag, TargetPath: targetMetaFlag}
			genMetaCmd.Run(nil, []string{})
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})
		It("Generate Mtad", func() {
			sourceMtadFlag = getTestPath("mtahtml5")
			ep = dir.MtaLocationParameters{SourcePath: sourceMtadFlag, TargetPath: targetMtadFlag}
			err := genMtadCmd.RunE(nil, []string{})
			if err != nil {
				fmt.Println(err)
			}
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})
		It("Generate Mtar", func() {
			sourceMtarFlag = getTestPath("mtahtml5")
			ep = dir.MtaLocationParameters{SourcePath: sourceMtarFlag, TargetPath: targetMtarFlag}
			genMetaCmd.Run(nil, []string{})
			genMtarCmd.Run(nil, []string{})
			Ω(getTestPath("result", "mtahtml5.mtar")).Should(BeAnExistingFile())
		})
	})

	var _ = Describe("Validate", func() {
		It("Invalid yaml path", func() {
			var str bytes.Buffer
			sourceValidateFlag = getTestPath("mta1")
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			validateCmd.Run(nil, []string{})

			Ω(str.String()).Should(ContainSubstring("Error reading the MTA file"))
		})
	})

	var _ = Describe("Generate Commands", func() {
		readFileContent := func(filename string) string {
			content, _ := ioutil.ReadFile(filename)
			contentString := string(content[:])
			contentString = strings.Replace(contentString, "\n", "", -1)
			contentString = strings.Replace(contentString, "\r", "", -1)
			return contentString
		}

		It("Generate Meta", func() {
			ep := dir.MtaLocationParameters{SourcePath: getTestPath("mtahtml5"), TargetPath: targetMetaFlag}
			generateMeta(&ep)
			mtadPath, _ := ep.GetMtadPath()
			Ω(readFileContent(mtadPath)).Should(Equal(readFileContent(getTestPath("golden", "mtad.yaml"))))
		})

		It("Generate Mtar", func() {
			ep := dir.MtaLocationParameters{SourcePath: getTestPath("mtahtml5"), TargetPath: targetMtarFlag}
			err := generateMeta(&ep)
			if err != nil {
				fmt.Println(err)
			}
			err = generateMtar(&ep)
			if err != nil {
				fmt.Println(err)
			}
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})

	var _ = Describe("Pack", func() {
		DescribeTable("Standard cases", func(projectPath string, validator func(ep *dir.MtaLocationParameters)) {
			sourcePackFlag = projectPath
			ep := dir.MtaLocationParameters{SourcePath: sourcePackFlag, TargetPath: targetPackFlag}
			pPackModuleFlag = "ui5app"
			packCmd.Run(nil, []string{})
			validator(&ep)
		},
			Entry("SanityTest",
				getTestPath("mtahtml5"),
				func(ep *dir.MtaLocationParameters) {
					Ω(ep.GetTargetModuleZipPath("ui5app")).Should(BeAnExistingFile())
				}),
			Entry("Wrong path to project",
				getTestPath("mtahtml6"),
				func(ep *dir.MtaLocationParameters) {
					Ω(ep.GetTargetModuleZipPath("ui5app")).ShouldNot(BeAnExistingFile())
				}),
		)
	})

	var _ = Describe("copyModuleArchive", func() {

		AfterEach(func() {
			os.Remove(getTestPath("mta", "node-js", "data.zip"))
			os.RemoveAll(getTestPath("mta", "mta"))
			dir.GetWorkingDirectory = dir.OsGetWd
		})

		It("Sanity", func() {
			lp := dir.MtaLocationParameters{SourcePath: getTestPath("mta")}
			dir.Archive(getTestPath("mta", "node-js"), getTestPath("mta", "node-js", "data.zip"))

			Ω(copyModuleArchive(&lp, "node-js", "node-js")).Should(Succeed())
			Ω(getTestPath("mta", "mta", "node-js", "data.zip")).Should(BeAnExistingFile())
		})

		var _ = DescribeTable("Invalid cases", func(modulePath string, mockWd bool, failOnCall int) {
			var lp dir.MtaLocationParameters
			var countCalls = 0
			if mockWd {
				lp = dir.MtaLocationParameters{}
				dir.GetWorkingDirectory = func() (string, error) {
					countCalls++
					if countCalls >= failOnCall {
						return "", errors.New("error")
					}
					return os.Getwd()
				}
			} else {
				lp = dir.MtaLocationParameters{SourcePath: getTestPath("mta")}
			}
			Ω(copyModuleArchive(&lp, modulePath, "node-js")).Should(HaveOccurred())
		},
			Entry("Invalid module name", "node-js1", false, -1),
			Entry("Get wd fails", "node-js", true, 1),
			Entry("Get wd fails", "node-js", true, 2))
	})

	var _ = Describe("Validation", func() {
		var _ = DescribeTable("getValidationMode", func(flag string, expectedValidateSchema, expectedValidateProject, expectedSuccess bool) {
			res1, res2, err := getValidationMode(flag)
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
			ep := dir.MtaLocationParameters{SourcePath: getTestPath(projectRelPath)}
			err := validateMtaYaml(&ep, validateSchema, validateProject)
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
			targetBModuleFlag = getTestPath("result")
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
			os.RemoveAll(targetBModuleFlag)
			builders.CommandsConfig = make([]byte, len(config))
			copy(builders.CommandsConfig, config)
		})

		It("Sanity", func() {
			ep := dir.MtaLocationParameters{SourcePath: getTestPath("mta"), TargetPath: targetBModuleFlag}
			Ω(buildModule(&ep, "node-js")).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})

		var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
			ep := dir.MtaLocationParameters{SourcePath: getTestPath(projectName), TargetPath: targetBModuleFlag, MtaFilename: mtaFilename}
			Ω(buildModule(&ep, moduleName)).Should(HaveOccurred())
			Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
		},
			Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
			Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
			Entry("Invalid module name", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
		)

		It("build Command", func() {
			pBuildModuleNameFlag = "node-js"
			sourceBModuleFlag = getTestPath("mta")
			ep := dir.MtaLocationParameters{SourcePath: sourceBModuleFlag, TargetPath: targetBModuleFlag}
			err := bModuleCmd.RunE(nil, []string{})
			if err != nil {
				fmt.Println(err)
			}
			Ω(ep.GetTargetModuleZipPath(pBuildModuleNameFlag)).Should(BeAnExistingFile())
		})
	})

	var _ = Describe("moduleCmd", func() {
		It("Sanity", func() {
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

			m := mta.MTA{}
			// parse mta yaml
			err := yaml.Unmarshal(mtaCF, &m)
			if err != nil {
				fmt.Println(err)
			}
			path, commands, err := moduleCmd(&m, "htmlapp")
			Ω(err).Should(BeNil())
			Ω(path).Should(Equal("app"))
			Ω(commands).Should(Equal([]string{"npm install", "grunt", "npm prune --production"}))
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
})

var _ = Describe("Process Dependencies", func() {
	It("Sanity", func() {
		Ω(processDependencies(&dir.MtaLocationParameters{SourcePath: getTestPath("mtahtml5"), MtaFilename: "mtaWithBuildParams.yaml"}, "ui5app")).Should(Succeed())
	})
	It("Invalid mta", func() {
		Ω(processDependencies(&dir.MtaLocationParameters{SourcePath: getTestPath("mtahtml5"), MtaFilename: "mta1.yaml"}, "ui5app")).Should(HaveOccurred())
	})
	It("Invalid module name", func() {
		Ω(processDependencies(&dir.MtaLocationParameters{SourcePath: getTestPath("mtahtml5")}, "xxx")).Should(HaveOccurred())
	})
})

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}
