package artifacts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/SAP/cloud-mta-build-tool/internal/fs"
	"github.com/SAP/cloud-mta/mta"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("manifest", func() {

	BeforeEach(func() {
		os.MkdirAll(getTestPath("result", "mta", "META-INF"), os.ModePerm)
	})

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})

	var _ = Describe("setManifestDesc", func() {
		It("Sanity", func() {
			os.Mkdir(filepath.Join(getTestPath("result", "mta"), "node-js"), os.ModePerm)
			os.Create(filepath.Join(getTestPath("result", "mta"), "node-js", "data.zip"))
			dirC, _ := ioutil.ReadDir(getTestPath("result", "mta"))
			for _, c := range dirC {
				fmt.Println(c.Name())
			}
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", "mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
	})

	var _ = Describe("genManifest", func() {
		It("Sanity", func() {
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			entries := []entry{
				{
					EntryName:   "node-js",
					EntryPath:   "node-js/data.zip",
					EntryType:   moduleEntry,
					ContentType: applicationZip,
				},
			}
			Ω(genManifest(loc.GetManifestPath(), entries)).Should(Succeed())
			actual := getFileContent(getTestPath("result", "mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
			Ω(actual).Should(Equal(golden))
		})
		It("Fails on wrong location", func() {
			loc := dir.Loc{}
			Ω(genManifest(loc.GetManifestPath(), []entry{})).Should(HaveOccurred())
		})
	})
	var _ = Describe("moduleDefined", func() {
		It("not defined", func() {
			Ω(moduleDefined("x", []string{"a"})).Should(BeFalse())
		})
		It("empty list", func() {
			Ω(moduleDefined("x", []string{})).Should(BeTrue())
		})
		It("defined", func() {
			Ω(moduleDefined("x", []string{"y", "x"})).Should(BeTrue())
		})

	})
})
