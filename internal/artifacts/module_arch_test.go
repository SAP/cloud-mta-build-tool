package artifacts

import (
	"os"

	"cloud-mta-build-tool/internal/builders"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
)

var _ = Describe("ModuleArch", func() {
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
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
				Ω(BuildModule(&ep, "node-js")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
				ep := dir.Loc{SourcePath: getTestPath(projectName), TargetPath: getTestPath("result"), MtaFilename: mtaFilename}
				Ω(BuildModule(&ep, moduleName)).Should(HaveOccurred())
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
			},
				Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
				Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
				Entry("Invalid module name", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
			)
		})

	})

	var _ = Describe("CopyModuleArchive", func() {

		It("Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta_with_zipped_module"), TargetPath: getTestPath("result")}
			Ω(CopyModuleArchive(&ep, "node-js", "node-js")).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})
		It("Invalid - no zip exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			Ω(CopyModuleArchive(&ep, "node-js", "node-js")).Should(HaveOccurred())
		})

		var _ = Describe("Invalid - Get Source failures", func() {

			AfterEach(func() {
				dir.GetWorkingDirectory = os.Getwd
			})

			DescribeTable("Failures", func(failOnCall int) {
				var callsCounter = 0
				wd, _ := os.Getwd()
				dir.GetWorkingDirectory = func() (string, error) {
					callsCounter++
					if callsCounter >= failOnCall {
						return "", errors.New("err")
					}
					return wd, nil
				}
				ep := dir.Loc{}
				Ω(CopyModuleArchive(&ep, "node-js", "node-js")).Should(HaveOccurred())
			},
				Entry("Fails on first call", 1),
				Entry("Fails on second call", 2))
		})

	})

})
