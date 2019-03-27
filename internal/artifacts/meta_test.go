package artifacts

import (
	"fmt"
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
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"), os.ModePerm)
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"), os.ModePerm)
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "CF", true, os.Getwd)).Should(Succeed())
			Ω(getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF", "MANIFEST.MF")).Should(BeAnExistingFile())
			Ω(getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF", "mtad.yaml")).Should(BeAnExistingFile())
		})

		It("Wrong location - fails on Working directory get", func() {
			Ω(ExecuteGenMeta("", "", "dev", "cf", true, func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Wrong platform", func() {
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"), os.ModePerm)
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"), os.ModePerm)
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "xx", true, os.Getwd)).Should(HaveOccurred())

		})
		It("generateMeta fails on wrong source path - parse mta fails", func() {
			Ω(ExecuteGenMeta(getTestPath("mtahtml6"), getResultPath(), "dev", "cf", true, os.Getwd)).Should(HaveOccurred())
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
			os.MkdirAll(getTestPath("result", ".testproject_mta_build_tmp", "app"), os.ModePerm)
			Ω(GenMetaInfo(&ep, &ep, ep.IsDeploymentDescriptor(), "cf",
				&m, []string{"htmlapp"}, true)).Should(Succeed())
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})

		It("Meta creation fails - fails on conversion by platform", func() {
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
			cfg := platform.PlatformConfig
			platform.PlatformConfig = []byte(`very bad config`)
			Ω(GenMetaInfo(&ep, &ep, ep.IsDeploymentDescriptor(), "cf",
				&m, []string{"htmlapp"}, true)).Should(HaveOccurred())
			platform.PlatformConfig = cfg
		})

		It("Fails on create file for manifest path", func() {
			loc := testLoc{ep}
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
			Ω(GenMetaInfo(&loc, &ep, ep.IsDeploymentDescriptor(), "cf",
				&m, []string{"htmlapp"}, true)).Should(HaveOccurred())
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
				Ω(GenMetaInfo(&ep, &ep, ep.IsDeploymentDescriptor(), "cf",
					&m, []string{"htmlapp"}, true)).Should(HaveOccurred())
			})
		})
	})

	var _ = Describe("Generate Commands", func() {

		readFileContent := func(ep dir.IMtaParser) *mta.MTA {
			mtaObj, _ := ep.ParseFile()
			return mtaObj
		}

		It("Generate Meta", func() {
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"), os.ModePerm)
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"), os.ModePerm)
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			generateMeta(&ep, &ep, &ep, false, "cf", true)
			Ω(readFileContent(&dir.Loc{SourcePath: getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF"), Descriptor: "dep"})).
				Should(Equal(readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})))
		})

		It("Generate Meta - with extension file", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"), os.ModePerm)
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"), os.ModePerm)
			Ω(generateMeta(&ep, &ep, &ep, false, "cf", true)).Should(Succeed())
			actual := readFileContent(&dir.Loc{SourcePath: getTestPath("result", ".mtahtml5_mta_build_tmp", "META-INF"), Descriptor: "dep"})
			golden := readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})
			Ω(actual).Should(Equal(golden))
		})

		It("Generate Meta - mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(),
				MtaFilename: "mtaNotExists.yaml"}
			Ω(generateMeta(&ep, &ep, &ep, false, "cf", true)).Should(HaveOccurred())
		})

		It("Generate Meta fails on platform parsing", func() {
			platformConfig := platform.PlatformConfig
			platform.PlatformConfig = []byte("wrong config")
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"), os.ModePerm)
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"), os.ModePerm)
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, &ep, false, "cf", true)).Should(HaveOccurred())
			platform.PlatformConfig = platformConfig
		})

		It("Generate Meta fails on mtad adaptation", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, &ep, false, "cf", true)).Should(HaveOccurred())
		})

		It("Generate Mtar", func() {
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "testapp"), os.ModePerm)
			os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp", "ui5app2"), os.ModePerm)
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			err := generateMeta(&ep, &ep, &ep, false, "cf", true)
			if err != nil {
				fmt.Println(err)
			}
			mtarPath, err := generateMtar(&ep, &ep, &ep, true, "")
			if err != nil {
				fmt.Println(err)
			}
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
	return loc.GetMtarDir(targetProvided)
}
