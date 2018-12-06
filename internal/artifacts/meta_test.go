package artifacts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

type testLoc struct {
	metaRes     string
	mtadRes     string
	manifestRes string
}

func (loc *testLoc) GetMetaPath() (string, error) {
	if loc.metaRes == "" {
		return "", errors.New("err")
	}
	return loc.metaRes, nil
}
func (loc *testLoc) GetMtadPath() (string, error) {
	if loc.mtadRes == "" {
		return "", errors.New("err")
	}
	return loc.mtadRes, nil
}

func (loc *testLoc) GetManifestPath() (string, error) {
	if loc.manifestRes == "" {
		return "", errors.New("err")
	}
	return loc.manifestRes, nil
}

var _ = Describe("Meta", func() {
	var _ = Describe("GenMetaInf", func() {
		wd, _ := os.Getwd()
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata", "testproject"), TargetPath: filepath.Join(wd, "testdata", "result")}
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
			Ω(GenMetaInfo(&ep, ep.IsDeploymentDescriptor(), "cf", &m, []string{"htmlapp"})).Should(Succeed())
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})

		DescribeTable("Invalid", func(loc *testLoc) {
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
			Ω(GenMetaInfo(loc, ep.IsDeploymentDescriptor(), "cf", &m, []string{"htmlapp"})).Should(HaveOccurred())
		},
			Entry("GenMtad fails", &testLoc{
				metaRes:     getTestPath("result", "META-INFO"),
				manifestRes: getTestPath("result", "META-INFO", "MANIFEST.MF"),
				mtadRes:     "",
			}),
			Entry("GetManifestPath fails", &testLoc{
				metaRes:     getTestPath("result", "META-INFO"),
				manifestRes: "",
				mtadRes:     getTestPath("result", "META-INFO", "mtad.yaml"),
			}))
	})

	var _ = Describe("Generate Commands", func() {

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		readFileContent := func(ep dir.ILoc) *mta.MTA {
			mtaObj, _ := ep.ParseFile()
			return mtaObj
		}

		It("Generate Meta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			GenerateMeta(&ep, "cf")
			Ω(readFileContent(&dir.Loc{SourcePath: getTestPath("result", "mtahtml5", "META-INF"), Descriptor: "dep"})).Should(Equal(readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})))
		})

		It("Generate Meta - with extension file", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			Ω(GenerateMeta(&ep, "cf")).Should(Succeed())
			actual := readFileContent(&dir.Loc{SourcePath: getTestPath("result", "mtahtml5", "META-INF"), Descriptor: "dep"})
			golden := readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})
			Ω(actual).Should(Equal(golden))
		})

		It("Generate Meta - mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaNotExists.yaml"}
			Ω(GenerateMeta(&ep, "cf")).Should(HaveOccurred())
		})

		It("Generate Mtar", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			err := GenerateMeta(&ep, "cf")
			if err != nil {
				fmt.Println(err)
			}
			err = GenerateMtar(&ep, &ep)
			if err != nil {
				fmt.Println(err)
			}
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})
})
