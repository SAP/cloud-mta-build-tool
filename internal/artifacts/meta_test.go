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
		os.RemoveAll(getResultPath())
	})

	var _ = Describe("ExecuteGenMeta", func() {

		It("Sanity", func() {
			createMtahtml5TmpFolder()
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "CF", os.Getwd)).Should(Succeed())
			Ω(getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF", "MANIFEST.MF")).Should(BeAnExistingFile())
			Ω(getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF", "mtad.yaml")).Should(BeAnExistingFile())
		})

		It("Fails on META-INF folder creation", func() {
			dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp"))
			file, _ := os.Create(getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF"))
			file.Close()
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "CF", os.Getwd)).Should(HaveOccurred())
		})

		It("Wrong location - fails on Working directory get", func() {
			Ω(ExecuteGenMeta("", "", "dev", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Wrong platform", func() {
			dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"))
			dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"))
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
			yaml.Unmarshal(mtaSingleModule, &m)
			dir.CreateDirIfNotExist(getTestPath("result", ".testproject_mta_build_tmp", "htmlapp"))
			os.Create(getTestPath("result", ".testproject_mta_build_tmp", "htmlapp", "data.zip"))
			dir.CreateDirIfNotExist(getTestPath("result", ".testproject_mta_build_tmp", "META-INF"))
			Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", &m)).Should(Succeed())
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})

		It("Meta creation fails - fails on conversion by platform", func() {
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
			dir.CreateDirIfNotExist(getTestPath("result", ".testproject_mta_build_tmp", "app"))
			dir.CreateDirIfNotExist(getTestPath("result", ".testproject_mta_build_tmp", "META-INF"))
			cfg := platform.PlatformConfig
			platform.PlatformConfig = []byte(`very bad config`)
			Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", &m)).Should(HaveOccurred())
			platform.PlatformConfig = cfg
		})

		It("Fails on create file for manifest path", func() {
			loc := testLoc{ep}
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
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
				os.RemoveAll(getResultPath())
			})

			It("Fails on get version", func() {
				m := mta.MTA{}
				yaml.Unmarshal(mtaSingleModule, &m)
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
			file, _ := os.Create(getTestPath("result", ".mtahtml5_mta_build_tmp", "xs-security.json"))
			file.Close()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, &ep, &ep, false, "cf")).Should(Succeed())
			Ω(readFileContent(&dir.Loc{SourcePath: getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF"), Descriptor: "dep"})).
				Should(Equal(readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})))
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
	dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"))
	dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"))
	dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app"))
	dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF"))
	createTmpFile(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app", "data.zip"))
	createTmpFile(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2", "data.zip"))
	createTmpFile(getTestPath("result", ".mtahtml5_mta_build_tmp", "xs-security.json"))
}

func createTmpFile(path string) {
	file, err := os.Create(path)
	file.Close()
	Ω(err).Should(Succeed())
}
