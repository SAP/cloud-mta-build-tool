package artifacts

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Meta", func() {

	AfterEach(func() {
		Ω(os.RemoveAll(getResultPath())).Should(Succeed())
	})

	var _ = Describe("ExecuteGenMeta", func() {

		It("Sanity", func() {
			createMtahtml5TmpFolder()
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "CF", os.Getwd)).Should(Succeed())
			Ω(getFullPathInTmpFolder("mtahtml5", "META-INF", "MANIFEST.MF")).Should(BeAnExistingFile())
			Ω(getFullPathInTmpFolder("mtahtml5", "META-INF", "mtad.yaml")).Should(BeAnExistingFile())
		})

		It("Fails on META-INF folder creation", func() {
			createDirInTempFolder("mtahtml5")
			createFileInTmpFolder("mtahtml5", "META-INF")
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "CF", os.Getwd)).Should(HaveOccurred())
		})

		It("Wrong location - fails on Working directory get", func() {
			Ω(ExecuteGenMeta("", "", "dev", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Wrong platform", func() {
			createDirInTempFolder("mtahtml5", "ui5app2")
			createDirInTempFolder("mtahtml5", "testapp")
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "xx", os.Getwd)).Should(HaveOccurred())

		})
		It("generateMeta fails on wrong source path - parse mta fails", func() {
			Ω(ExecuteGenMeta(getTestPath("mtahtml6"), getResultPath(), "dev", "cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("GenMetaInfo", func() {
		ep := dir.Loc{SourcePath: getTestPath("testproject"), TargetPath: getResultPath()}
		var mtaSingleModule = []byte(`
_schema-version: "2.0.0"
ID: mta_proj
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
`)

		It("Sanity", func() {
			m := mta.MTA{}
			Ω(yaml.Unmarshal(mtaSingleModule, &m)).Should(Succeed())
			createDirInTempFolder("testproject", "htmlapp")
			createFileInTmpFolder("testproject", "htmlapp", "data.zip")
			createDirInTempFolder("testproject", "META-INF")
			Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", &m)).Should(Succeed())
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})

		It("Meta creation fails - fails on conversion by platform", func() {
			m := mta.MTA{}
			Ω(yaml.Unmarshal(mtaSingleModule, &m)).Should(Succeed())
			createDirInTempFolder("testproject", "app")
			createDirInTempFolder("testproject", "META-INF")
			cfg := platform.PlatformConfig
			platform.PlatformConfig = []byte(`very bad config`)
			Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", &m)).Should(HaveOccurred())
			platform.PlatformConfig = cfg
		})

		It("Fails on create file for manifest path", func() {
			loc := testLoc{ep}
			m := mta.MTA{}
			Ω(yaml.Unmarshal(mtaSingleModule, &m)).Should(Succeed())
			Ω(genMetaInfo(&loc, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", &m)).Should(HaveOccurred())
		})

		var _ = Describe("Fails on setManifestDesc", func() {
			var config []byte

			BeforeEach(func() {
				config = make([]byte, len(version.VersionConfig))
				copy(config, version.VersionConfig)
				// Simplified commands configuration (performance purposes). removed "npm prune --production"
				version.VersionConfig = []byte(`
cli_version:["x"]
`)
			})

			AfterEach(func() {
				version.VersionConfig = make([]byte, len(config))
				copy(version.VersionConfig, config)
				Ω(os.RemoveAll(getResultPath())).Should(Succeed())
			})

			It("Fails on get version", func() {
				m := mta.MTA{}
				Ω(yaml.Unmarshal(mtaSingleModule, &m)).Should(Succeed())
				Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", &m)).Should(HaveOccurred())
			})
		})
	})

	var _ = Describe("Generate Commands", func() {

		readFileContent := func(ep dir.IMtaParser) *mta.MTA {
			mtaObj, _ := ep.ParseFile()
			return mtaObj
		}

		It("Generate Meta", func() {
			createMtahtml5TmpFolder()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, &ep, &ep, false, "cf")).Should(Succeed())
			Ω(readFileContent(&dir.Loc{SourcePath: getFullPathInTmpFolder("mtahtml5", "META-INF"), Descriptor: "dep"})).
				Should(Equal(readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})))
		})
		It("Generate Mtad - fails on missing paths in temp folder", func() {
			createDirInTempFolder("mtahtml5", "META-INF")
			createFileInTmpFolder("mtahtml5", "xs-security.json")
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			mtaStr, err := ep.ParseFile()
			Ω(err).Should(Succeed())
			for i := range mtaStr.Modules {
				mtaStr.Modules[i].Path = ""
			}
			err = genMetaInfo(&ep, &ep, &ep, false, "neo", mtaStr)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`failed to clean the build parameters from the "ui5app" module`))
		})

		It("Generate Meta - mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(),
				MtaFilename: "mtaNotExists.yaml"}
			err := generateMeta(&ep, &ep, &ep, &ep, false, "cf")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`failed to read the "%s"`, ep.GetMtaYamlPath()))
		})

		It("Generate Meta fails on platform parsing", func() {
			platformConfig := platform.PlatformConfig
			platform.PlatformConfig = []byte("wrong config")
			createMtahtml5TmpFolder()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			err := generateMeta(&ep, &ep, &ep, &ep, false, "cf")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`generation of the MTAD file failed when converting types according to the "cf" platform: unmarshalling of the platforms failed`))
			platform.PlatformConfig = platformConfig
		})

		It("Generate Meta fails on mtad adaptation", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, &ep, &ep, false, "cf")).Should(HaveOccurred())
		})

		It("Generate Mtar", func() {
			createMtahtml5TmpFolder()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			err := generateMeta(&ep, &ep, &ep, &ep, false, "cf")
			Ω(err).Should(Succeed())
			mtarPath, err := generateMtar(&ep, &ep, &ep, true, "")
			Ω(err).Should(Succeed())
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})
})

type testLoc struct {
	loc dir.Loc
}

func (loc *testLoc) GetMetaPath() string {
	return loc.loc.GetMetaPath()
}

func (loc *testLoc) GetMtadPath() string {
	return loc.loc.GetMtadPath()
}

func (loc *testLoc) GetManifestPath() string {
	return filepath.Join(loc.loc.GetManifestPath(), "folderNotExists", "MANIFEST.MF")
}

func (loc *testLoc) GetMtarDir(targetProvided bool) string {
	return loc.loc.GetMtarDir(targetProvided)
}

func (loc *testLoc) GetSourceModuleDir(modulePath string) string {
	return loc.loc.GetSourceModuleDir(modulePath)
}

func createMtahtml5TmpFolder() {
	createDirInTempFolder("mtahtml5", "ui5app2")
	createDirInTempFolder("mtahtml5", "ui5app")
	createDirInTempFolder("mtahtml5", "META-INF")
	createFileInTmpFolder("mtahtml5", "ui5app", "data.zip")
	createFileInTmpFolder("mtahtml5", "ui5app2", "data.zip")
	createFileInTmpFolder("mtahtml5", "xs-security.json")
}
