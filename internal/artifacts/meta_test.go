package artifacts

import (
	"fmt"
	"os"

	"cloud-mta-build-tool/internal/platform"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("Meta", func() {

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})

	var _ = Describe("ExecuteGenMeta", func() {

		It("Sanity", func() {
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getTestPath("result"), "dev", "cf", os.Getwd)).Should(Succeed())
			Ω(getTestPath("result", "mtahtml5", "META-INF", "MANIFEST.MF")).Should(BeAnExistingFile())
			Ω(getTestPath("result", "mtahtml5", "META-INF", "mtad.yaml")).Should(BeAnExistingFile())
		})

		It("Wrong location - fails on Working directory get", func() {
			Ω(ExecuteGenMeta("", "", "dev", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("generateMeta fails on wrong source path - parse mta fails", func() {
			Ω(ExecuteGenMeta(getTestPath("mtahtml6"), getTestPath("result"), "dev", "cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("GenMetaInf", func() {
		ep := dir.Loc{SourcePath: getTestPath("testproject"), TargetPath: getTestPath("result")}
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
	})

	var _ = Describe("Generate Commands", func() {

		readFileContent := func(ep dir.IMtaParser) *mta.MTA {
			mtaObj, _ := ep.ParseFile()
			return mtaObj
		}

		It("Generate Meta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			generateMeta(&ep, &ep, false, "cf")
			Ω(readFileContent(&dir.Loc{SourcePath: getTestPath("result", "mtahtml5", "META-INF"), Descriptor: "dep"})).Should(Equal(readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})))
		})

		It("Generate Meta - with extension file", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			Ω(generateMeta(&ep, &ep, false, "cf")).Should(Succeed())
			actual := readFileContent(&dir.Loc{SourcePath: getTestPath("result", "mtahtml5", "META-INF"), Descriptor: "dep"})
			golden := readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})
			Ω(actual).Should(Equal(golden))
		})

		It("Generate Meta - mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaNotExists.yaml"}
			Ω(generateMeta(&ep, &ep, false, "cf")).Should(HaveOccurred())
		})

		It("Generate Meta fails on platform parsing", func() {
			platformConfig := platform.PlatformConfig
			platform.PlatformConfig = []byte("wrong")
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			Ω(generateMeta(&ep, &ep, false, "cf")).Should(HaveOccurred())
			platform.PlatformConfig = platformConfig
		})

		It("Generate Mtar", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			err := generateMeta(&ep, &ep, false, "cf")
			if err != nil {
				fmt.Println(err)
			}
			err = generateMtar(&ep, &ep)
			if err != nil {
				fmt.Println(err)
			}
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})
})
