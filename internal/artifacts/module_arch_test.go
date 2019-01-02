package artifacts

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/commands"
	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("ModuleArch", func() {

	var config []byte

	BeforeEach(func() {
		config = make([]byte, len(commands.CommandsConfig))
		copy(config, commands.CommandsConfig)
		// Simplified commands configuration (performance purposes). removed "npm prune --production"
		commands.CommandsConfig = []byte(`
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
		commands.CommandsConfig = make([]byte, len(config))
		copy(commands.CommandsConfig, config)
		os.RemoveAll(getTestPath("result"))
	})

	m := mta.Module{
		Name: "node-js",
		Path: "node-js",
	}

	var _ = Describe("ExecuteBuild", func() {

		It("Sanity", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getTestPath("result"), "dev", "node-js", "cf", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())

		})

		It("Fails on location initialization", func() {
			Ω(ExecuteBuild("", "", "dev", "ui5app", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getTestPath("result"), "dev", "ui5app", "cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("ExecutePack", func() {

		It("Sanity", func() {
			Ω(ExecutePack(getTestPath("mta"), getTestPath("result"), "dev", "node-js",
				"cf", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})

		It("Fails on location initialization", func() {
			Ω(ExecutePack("", "", "dev", "ui5app", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecutePack(getTestPath("mta"), getTestPath("result"), "dev", "ui5appx",
				"cf", os.Getwd)).Should(HaveOccurred())
		})

		It("Target folder exists as file", func() {
			os.MkdirAll(getTestPath("result", "mta"), os.ModePerm)
			createFile("result", "mta", "node-js")
			Ω(ExecutePack(getTestPath("mta"), getTestPath("result"), "dev", "node-js",
				"cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("Pack", func() {
		var _ = Describe("Sanity", func() {

			It("Deployment descriptor - Copy only", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getTestPath("result"),
					Descriptor: "dep",
				}
				Ω(packModule(&ep, true, &m, "node-js", "cf")).Should(Succeed())
				Ω(getTestPath("result", "mta_with_zipped_module", "node-js", "data.zip")).Should(BeAnExistingFile())
			})

			//ep.GetTargetModuleDir(moduleName)
			It("Wrong source", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_unknown"),
					TargetPath: getTestPath("result"),
					Descriptor: "dev",
				}
				Ω(packModule(&ep, false, &m, "node-js", "cf")).Should(HaveOccurred())
			})
			It("Target directory exists as a file", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getTestPath("result"),
					Descriptor: "dev",
				}
				os.MkdirAll(filepath.Join(ep.GetTarget(), "mta_with_zipped_module"), os.ModePerm)
				createFile("result", "mta_with_zipped_module", "node-js")
				Ω(packModule(&ep, false, &m, "node-js", "cf")).Should(HaveOccurred())
			})
		})

		It("No platforms - no pack", func() {
			ep := dir.Loc{
				SourcePath: getTestPath("mta_with_zipped_module"),
				TargetPath: getTestPath("result"),
				Descriptor: "dep",
			}
			mNoPlatforms := mta.Module{
				Name: "node-js",
				Path: "node-js",
				BuildParams: map[string]interface{}{
					buildops.SupportedPlatformsParam: []string{},
				},
			}
			Ω(packModule(&ep, false, &mNoPlatforms, "node-js", "cf")).Should(Succeed())
			Ω(getTestPath("result", "mta_with_zipped_module", "node-js", "data.zip")).
				ShouldNot(BeAnExistingFile())
		})

	})

	var _ = Describe("Build", func() {

		var _ = Describe("build Module", func() {

			var config []byte

			BeforeEach(func() {
				config = make([]byte, len(commands.CommandsConfig))
				copy(config, commands.CommandsConfig)
				// Simplified commands configuration (performance purposes). removed "npm prune --production"
				commands.CommandsConfig = []byte(`
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

			It("Sanity", func() {
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
				Ω(buildModule(&ep, &ep, false, "node-js", "cf")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			It("Commands fail", func() {
				commands.CommandsConfig = []byte(`
builders:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  type:
    - command: go test exec_unknownTest.go
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  type:
    - command: go test exec_unknownTest.go
`)

				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
				Ω(buildModule(&ep, &ep, false, "node-js", "cf")).Should(HaveOccurred())
			})

			It("Target folder exists as a file - dev", func() {
				os.MkdirAll(getTestPath("result", "mta"), os.ModePerm)
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
				createFile("result", "mta", "node-js")
				Ω(buildModule(&ep, &ep, false, "node-js", "cf")).Should(HaveOccurred())
			})

			It("Target folder exists as a file - dep", func() {
				os.MkdirAll(getTestPath("result", "mta"), os.ModePerm)
				ep := dir.Loc{
					SourcePath:  getTestPath("mta"),
					TargetPath:  getTestPath("result"),
					Descriptor:  "dep",
					MtaFilename: "mta.yaml",
				}
				createFile("result", "mta", "node-js")
				Ω(buildModule(&ep, &ep, true, "node-js", "cf")).Should(HaveOccurred())
			})

			It("Deployment Descriptor", func() {
				ep := dir.Loc{
					SourcePath:  getTestPath("mta_with_zipped_module"),
					TargetPath:  getTestPath("result"),
					MtaFilename: "mta.yaml",
					Descriptor:  "dep"}
				Ω(buildModule(&ep, &ep, true, "node-js", "cf")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
				ep := dir.Loc{SourcePath: getTestPath(projectName), TargetPath: getTestPath("result"), MtaFilename: mtaFilename}
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
				Ω(buildModule(&ep, &ep, false, moduleName, "cf")).Should(HaveOccurred())
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
			},
				Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
				Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
				Entry("Invalid module name wrong build params", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
			)
		})
	})

	var _ = Describe("copyModuleArchive", func() {

		It("Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta_with_zipped_module"), TargetPath: getTestPath("result")}
			Ω(copyModuleArchive(&ep, "node-js", "node-js")).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})
		It("Invalid - no zip exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			Ω(copyModuleArchive(&ep, "node-js", "node-js")).Should(HaveOccurred())
		})
		It("Target directory exists as file", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta_with_zipped_module"), TargetPath: getTestPath("result")}
			os.MkdirAll(getTestPath("result", "mta_with_zipped_module"), os.ModePerm)
			createFile("result", "mta_with_zipped_module", "node-js")
			Ω(copyModuleArchive(&ep, "node-js", "node-js")).Should(HaveOccurred())
		})
	})
})

func createFile(path ...string) {
	file, err := os.Create(getTestPath(path...))
	Ω(err).Should(Succeed())
	file.Close()
}
