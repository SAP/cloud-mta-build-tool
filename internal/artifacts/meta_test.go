package artifacts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("Meta", func() {
	var _ = Describe("GenMetaInf", func() {
		wd, _ := os.Getwd()
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata", "testproject"), TargetPath: filepath.Join(wd, "testdata", "result")}
		epCurrent := dir.Loc{}
		var mtaSingleModule = []byte(`
_schema-version: "2.0.0"
ID: mta_proj
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
`)

		AfterEach(func() {
			targetDir, _ := ep.GetTarget()
			os.RemoveAll(targetDir)
			dir.GetWorkingDirectory = dir.OsGetWd
		})

		It("Sanity", func() {
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
			Ω(GenMetaInfo(&ep, "cf", &m, []string{"htmlapp"}, func(mtaStr *mta.MTA, platform string) {})).Should(Succeed())
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})

		It("Invalid", func() {
			dir.GetWorkingDirectory = func() (string, error) {
				return "", errors.New("err")
			}
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
			Ω(GenMetaInfo(&epCurrent, "cf", &m, []string{"htmlapp"}, func(mtaStr *mta.MTA, platform string) {})).Should(HaveOccurred())
		})
	})

	var _ = Describe("Generate Commands", func() {

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		readFileContent := func(ep *dir.Loc) *mta.MTA {
			mtaObj, _ := dir.ParseFile(ep)
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
			err = GenerateMtar(&ep)
			if err != nil {
				fmt.Println(err)
			}
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})
})
