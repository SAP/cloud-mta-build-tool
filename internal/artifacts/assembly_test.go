package artifacts

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/contenttype"
	"github.com/SAP/cloud-mta-build-tool/internal/fs"
)

var _ = Describe("Assembly", func() {

	BeforeEach(func() {
		os.Mkdir(getResultPath(), os.ModePerm)
	})

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})

	var _ = Describe("Assembly", func() {

		It("Sanity", func() {
			Ω(Assembly(getTestPath("assembly-sample"), getTestPath("result"), os.Getwd)).Should(Succeed())
			Ω(getTestPath("result", "com.sap.xs2.samples.javahelloworld.mtar")).Should(BeAnExistingFile())
		})
		It("Fails on getting working directory", func() {
			Ω(Assembly("", getTestPath("result"), func() (string, error) {
				return "", errors.New("error")
			})).Should(HaveOccurred())
		})
		It("Wrong source path - fails on parsing the .mtad file", func() {
			Ω(Assembly(getTestPath("assembly-sample1"), getTestPath("result"), os.Getwd)).Should(HaveOccurred())
		})
		It("Temporary folder exists as file", func() {
			file, err := os.Create(getTestPath("result", "assembly-sample"))
			Ω(err).Should(Succeed())
			file.Close()
			Ω(Assembly(getTestPath("assembly-sample"), getTestPath("result"), os.Getwd)).Should(HaveOccurred())
		})
		It("Entries missing", func() {
			Ω(Assembly(getTestPath("assembly-sample_broken"), getTestPath("result"), os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("getAssembledEntries", func() {
		It("Sanity", func() {
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), Descriptor: "dep"}
			mta, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			entries, err := getAssembledEntries(&loc, mta)
			Ω(err).Should(Succeed())
			Ω(len(entries)).Should(Equal(3))
		})
		var _ = DescribeTable("Mtad with broken paths", func(filename string) {
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), Descriptor: "dep", MtaFilename: filename}
			mta, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			_, err = getAssembledEntries(&loc, mta)
			Ω(err).Should(HaveOccurred())
		},
			Entry("Broken path of module", "mtadBrokenPath.yaml"),
			Entry("Broken path of requires", "mtadBrokenPathInRequires.yaml"),
			Entry("Broken path of resources", "mtadBrokenPathInResources.yaml"))
	})
	It("Wrong content types", func() {
		config := contenttype.ContentTypeConfig
		contenttype.ContentTypeConfig = []byte("Wrong content type config")
		loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), Descriptor: "dep"}
		mta, err := loc.ParseFile()
		Ω(err).Should(Succeed())
		_, err = getAssembledEntries(&loc, mta)
		Ω(err).Should(HaveOccurred())
		contenttype.ContentTypeConfig = config
	})
	It("Missing content types", func() {
		config := contenttype.ContentTypeConfig
		contenttype.ContentTypeConfig = []byte(`
content-types:
- extension: .war
  content-type: "application/zip"
- extension: .yaml
  content-type: "text/plain"		
`)
		loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), Descriptor: "dep"}
		mta, err := loc.ParseFile()
		Ω(err).Should(Succeed())
		_, err = getAssembledEntries(&loc, mta)
		Ω(err).Should(HaveOccurred())
		contenttype.ContentTypeConfig = config
	})

})
