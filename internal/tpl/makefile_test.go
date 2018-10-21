package tpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		makeFileName     = "MakeFileTest.mta"
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
		makeFileFullPath = func() string {
			path, _ := dir.GetFullPath("testdata", makeFileName)
			return path
		}()
		expectedMakeFileContent = getMakeFileContent(expectedMakePath)
	)

	var _ = Describe("MakeFile Generation", func() {
		AfterEach(func() {
			e := os.Remove(makeFileFullPath)
			if e != nil {
				fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaa " + makeFileFullPath)
			}
		})

		It("createMakeFile testing", func() {
			makeFilePath, _ := dir.GetFullPath("testdata")
			file, _ := createMakeFile(makeFilePath, makeFileName)
			Ω(file).ShouldNot(BeNil())
			file.Close()
			Ω(makeFilePath).Should(BeAnExistingFile())
			Ω(createMakeFile(makeFilePath, makeFileName)).Should(BeNil())
		})
		It("Sanity", func() {
			Ω(makeFile(makeFileName, tpl)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileContent))
		})
		It("Make testing with wrong mode", func() {
			Ω(Make("wrongMode")).Should(HaveOccurred())
		})
	})

	var _ = DescribeTable("Makefile Generation Failed", func(testPath, testTemplate string) {

		Ω(makeFile(makeFileName, tplCfg{relPath: testPath, tplName: testTemplate, pre: basePreVerbose, post: basePostVerbose})).Should(HaveOccurred())
	},
		Entry("Wrong Template", "testdata", filepath.Join("testdata", "WrongMakeTmpl.txt")),
		Entry("Yaml not exists", "testdata1", "make_default.txt"),
	)

	var _ = DescribeTable("String in slice search", func(s string, slice []string, expected bool) {
		Ω(stringInSlice(s, slice)).Should(Equal(expected))
	},
		Entry("positive test", "test1", []string{"test1", "foo"}, true),
		Entry("negative test", "test1", []string{"--test", "foo"}, false),
	)

	var _ = Describe("Make mode tests", func() {
		DescribeTable("Positive", func(mode string, tpl tplCfg) {
			Ω(makeMode(mode)).Should(Equal(tpl))
		},
			Entry("Default mode", "", tplCfg{tplName: makeDefaultTpl, pre: basePreDefault, post: basePostDefault}),
			Entry("Verbose mode", "verbose", tplCfg{tplName: makeVerboseTpl, pre: basePreVerbose, post: basePostVerbose}),
		)
		It("unknown mode", func() {
			_, err := makeMode("test")
			Ω(err).Should(MatchError("command is not supported"))
		})
	},
	)

})
