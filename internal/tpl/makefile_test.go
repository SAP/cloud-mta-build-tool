package tpl

import (
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})

func removeSpecialSymbols(b []byte) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	s := string(b)
	s = strings.Replace(s, "0xd, ", "", -1)
	s = reg.ReplaceAllString(s, "")
	return s
}

func getMakeFileContent(filePath string) string {
	expected, _ := ioutil.ReadFile(filePath)
	return removeSpecialSymbols(expected)
}

var _ = Describe("Makefile", func() {

	var (
		tpl              = tplCfg{tplName: "make_verbose.txt", relPath: "testdata", pre: basePreVerbose, post: basePostVerbose}
		makeFileName     = "MakeFileTest"
		expectedMakePath = func() string {
			var filename string
			switch runtime.GOOS {
			case "linux":
				filename = "ExpectedMakeFileLinux"
			case "darwin":
				filename = "ExpectedMakeFileMac"
			default:
				filename = "ExpectedMakeFileWindows"
			}
			path, _ := dir.GetFullPath("testdata", filename)
			return path
		}()
		makeFilePath = func() string {
			path, _ := dir.GetFullPath("testdata", makeFileName)
			return path
		}()
		makeFileExtendedPath    = makeFilePath + ".mta"
		expectedMakeFileContent = getMakeFileContent(expectedMakePath)
	)

	var _ = Describe("MakeFile Generation", func() {
		assertMakeFile := func(expectedMakeFilePath string) {
			Ω(makeFile(makeFileName, tpl)).Should(Succeed())
			Ω(expectedMakeFilePath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(expectedMakeFilePath)).Should(Equal(expectedMakeFileContent))
		}
		AfterEach(func() {
			os.Remove(makeFilePath)
			os.Remove(makeFileExtendedPath)
		})

		It("Sanity", func() {
			assertMakeFile(makeFilePath)
			Ω(makeFileExtendedPath).ShouldNot(BeAnExistingFile())
			assertMakeFile(makeFileExtendedPath)
		})
	})

	var _ = DescribeTable("Makefile Generation Failed", func() {

		//Ω(makeFile(makeFileName, tplCfg{tplName: filepath.Join("testdata", tplFilename)})).Should(HaveOccurred())
	},
	//Entry("Wrong Template", "WrongMakeTmpl.txt"),
	//Entry("Empty Template", "emptyMakeTmpl.txt"),
	)

})
