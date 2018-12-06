package artifacts

import (
	"os"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/mta"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
)

type testPackLoc struct {
	targetDir string
	sourceDir string
	zipPath   string
}

func (loc *testPackLoc) GetSourceModuleDir(modulePath string) (string, error) {
	if loc.sourceDir == "" {
		return "", errors.New("err")
	}
	return loc.sourceDir, nil
}
func (loc *testPackLoc) GetTargetModuleDir(moduleName string) (string, error) {
	if loc.targetDir == "" {
		return "", errors.New("err")
	}
	return loc.targetDir, nil
}
func (loc *testPackLoc) GetTargetModuleZipPath(moduleName string) (string, error) {
	if loc.zipPath == "" {
		return "", errors.New("err")
	}
	return loc.zipPath, nil
}

var _ = Describe("ModuleArch", func() {

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	m := mta.Module{
		Name: "node-js",
		Path: "node-js",
	}
	var _ = Describe("Pack", func() {
		It("Deployment descriptor - Copy only", func() {
			ep := dir.Loc{
				SourcePath: getTestPath("mta_with_zipped_module"),
				TargetPath: getTestPath("result"),
				Descriptor: "dep",
			}
			Ω(PackModule(&ep, true, &m, "node-js")).Should(Succeed())
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
			Ω(PackModule(&ep, false, &mNoPlatforms, "node-js")).Should(Succeed())
			Ω(getTestPath("result", "mta_with_zipped_module", "node-js", "data.zip")).ShouldNot(BeAnExistingFile())
		})

		var _ = DescribeTable("Failures", func(loc *testPackLoc) {
			Ω(PackModule(loc, false, &m, "node-js")).Should(HaveOccurred())
		},
			Entry("GetTargetModuleDir fails", &testPackLoc{
				targetDir: "",
				sourceDir: getTestPath("mta"),
				zipPath:   "",
			}),
			Entry("GetBuildResultsPath fails", &testPackLoc{
				targetDir: getTestPath("result"),
				sourceDir: "",
				zipPath:   "",
			}))
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
				Ω(BuildModule(&ep, "node-js")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			It("Deployment Descriptor", func() {
				ep := dir.Loc{
					SourcePath:  getTestPath("mta_with_zipped_module"),
					TargetPath:  getTestPath("result"),
					MtaFilename: "mta.yaml",
					Descriptor:  "dep"}
				Ω(BuildModule(&ep, "node-js")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
				ep := dir.Loc{SourcePath: getTestPath(projectName), TargetPath: getTestPath("result"), MtaFilename: mtaFilename}
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
				Ω(BuildModule(&ep, moduleName)).Should(HaveOccurred())
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
			},
				Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
				Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
				Entry("Invalid module name wrong build params", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
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

			DescribeTable("Failures", func(loc *testPackLoc) {
				Ω(CopyModuleArchive(loc, "node-js", "node-js")).Should(HaveOccurred())
			},
				Entry("Fails on first call", &testPackLoc{
					sourceDir: "",
				}),
				Entry("Fails on second call", &testPackLoc{
					sourceDir: getTestPath(),
				}))
		})
	})
})
