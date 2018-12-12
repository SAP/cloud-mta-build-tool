package artifacts

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("ModuleArch", func() {

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
		builders.CommandsConfig = make([]byte, len(config))
		copy(builders.CommandsConfig, config)
		os.RemoveAll(getTestPath("result"))
	})

	m := mta.Module{
		Name: "node-js",
		Path: "node-js",
	}

	var _ = Describe("ExecuteBuild", func() {

		It("Sanity", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getTestPath("result"), "dev", "node-js", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())

		})

		It("Fails on location initialization", func() {
			Ω(ExecuteBuild("", "", "dev", "ui5app", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getTestPath("result"), "dev", "ui5app", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("ExecutePack", func() {
		It("Sanity", func() {
			Ω(ExecutePack(getTestPath("mta"), getTestPath("result"), "dev", "node-js", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})

		It("Fails on location initialization", func() {
			Ω(ExecutePack("", "", "dev", "ui5app", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecutePack(getTestPath("mta"), getTestPath("result"), "dev", "ui5appx", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("Pack", func() {
		It("Deployment descriptor - Copy only", func() {
			ep := dir.Loc{
				SourcePath: getTestPath("mta_with_zipped_module"),
				TargetPath: getTestPath("result"),
				Descriptor: "dep",
			}
			Ω(packModule(&ep, true, &m, "node-js")).Should(Succeed())
			Ω(getTestPath("result", "mta_with_zipped_module", "node-js", "data.zip")).Should(BeAnExistingFile())
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
			Ω(packModule(&ep, false, &mNoPlatforms, "node-js")).Should(Succeed())
			Ω(getTestPath("result", "mta_with_zipped_module", "node-js", "data.zip")).ShouldNot(BeAnExistingFile())
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

			It("Sanity", func() {
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
				Ω(buildModule(&ep, &ep, false, "node-js")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			It("Deployment Descriptor", func() {
				ep := dir.Loc{
					SourcePath:  getTestPath("mta_with_zipped_module"),
					TargetPath:  getTestPath("result"),
					MtaFilename: "mta.yaml",
					Descriptor:  "dep"}
				Ω(buildModule(&ep, &ep, true, "node-js")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
				ep := dir.Loc{SourcePath: getTestPath(projectName), TargetPath: getTestPath("result"), MtaFilename: mtaFilename}
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
				Ω(buildModule(&ep, &ep, false, moduleName)).Should(HaveOccurred())
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
	})
})
