package artifacts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
			Ω(GenMetaInfo(&ep, &m, []string{"htmlapp"}, func(mtaStr *mta.MTA) {})).Should(Succeed())
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})

		It("Invalid", func() {
			dir.GetWorkingDirectory = func() (string, error) {
				return "", errors.New("err")
			}
			m := mta.MTA{}
			yaml.Unmarshal(mtaSingleModule, &m)
			Ω(GenMetaInfo(&epCurrent, &m, []string{"htmlapp"}, func(mtaStr *mta.MTA) {})).Should(HaveOccurred())
		})
	})

	var _ = Describe("Generate Commands", func() {

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		readFileContent := func(filename string) string {
			content, _ := ioutil.ReadFile(filename)
			contentString := string(content[:])
			contentString = strings.Replace(contentString, "\n", "", -1)
			contentString = strings.Replace(contentString, "\r", "", -1)
			return contentString
		}

		It("Generate Meta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			GenerateMeta(&ep)
			mtadPath, _ := ep.GetMtadPath()
			Ω(readFileContent(mtadPath)).Should(Equal(readFileContent(getTestPath("golden", "mtad.yaml"))))
		})

		It("Generate Mtar", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			err := GenerateMeta(&ep)
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
