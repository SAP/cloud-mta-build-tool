package artifacts

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, mtaObj.Modules, []string{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", "mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
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
