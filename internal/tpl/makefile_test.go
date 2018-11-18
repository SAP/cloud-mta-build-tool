package tpl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/version"
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
		tpl              = tplCfg{tplName: makeVerboseTpl, relPath: "", pre: basePreVerbose, post: basePostVerbose, depDesc: "dev"}
		tplDep           = tplCfg{tplName: makeVerboseDepTpl, relPath: "", pre: basePreVerbose, post: basePostVerbose, depDesc: "dep"}
		makeFileName     = "MakeFileTest.mta"
		wd, _            = os.Getwd()
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
			return filepath.Join(wd, "testdata", filename)
		}()
		expectedMakeDepPath = func() string {
			var filename string
			switch runtime.GOOS {
			case "linux":
				filename = "ExpectedMakeFileDepLinux"
			case "darwin":
				filename = "ExpectedMakeFileDepMac"
			default:
				filename = "ExpectedMakeFileDepWindows"
			}
			return filepath.Join(wd, "testdata", filename)
		}()
		makeFileFullPath = func() string {
			return filepath.Join(wd, "testdata", makeFileName)
		}()
		expectedMakeFileContent    = getMakeFileContent(expectedMakePath)
		expectedMakeFileDepContent = getMakeFileContent(expectedMakeDepPath)
	)

	var _ = Describe("MakeFile Generation", func() {
		BeforeEach(func() {
			version.VersionConfig = []byte(`
cli_version: 0.0.0
makefile_version: 0.0.0
`)
		})
		AfterEach(func() {
			os.Remove(makeFileFullPath)
		})

		It("createMakeFile testing", func() {
			makeFilePath := filepath.Join(wd, "testdata")
			file, _ := createMakeFile(makeFilePath, makeFileName)
			Ω(file).ShouldNot(BeNil())
			file.Close()
			Ω(makeFilePath).Should(BeAnExistingFile())
			Ω(createMakeFile(makeFilePath, makeFileName)).Should(BeNil())
		})
		It("Sanity - Dev", func() {
			Ω(makeFile(&dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata"), Descriptor: "dev"}, makeFileName, &tpl)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileContent))
		})
		It("Sanity - Dep", func() {
			Ω(makeFile(&dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata"), Descriptor: "dep"}, makeFileName, &tplDep)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileDepContent))
		})
		It("Make testing with wrong mta yaml file", func() {
			Ω(Make(&dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "xxx.yaml"}, "")).Should(HaveOccurred())
		})
		It("Make testing with wrong mode", func() {
			Ω(Make(&dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata")}, "wrongMode")).Should(HaveOccurred())
		})
	})

	var _ = DescribeTable("Makefile Generation Failed", func(testPath, testTemplate string) {
		ep := dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata")}
		Ω(makeFile(&ep, makeFileName, &tplCfg{relPath: testPath, tplName: testTemplate, pre: basePreVerbose, post: basePostVerbose})).Should(HaveOccurred())
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
		DescribeTable("Positive", func(mode string, tpl tplCfg, isDep bool) {
			Ω(getTplCfg(mode, isDep)).Should(Equal(tpl))
		},
			Entry("Default mode Dev", "", tplCfg{tplName: makeDefaultTpl, pre: basePreDefault, post: basePostDefault}, false),
			Entry("Default mode Dep", "", tplCfg{tplName: makeDefaultTpl, pre: basePreDefault, post: basePostDefault}, true),
			Entry("Verbose mode Dev", "verbose", tplCfg{tplName: makeVerboseTpl, pre: basePreVerbose, post: basePostVerbose}, false),
			Entry("Verbose mode Dep", "verbose", tplCfg{tplName: makeVerboseDepTpl, pre: basePreVerbose, post: basePostVerbose}, true),
		)
		It("unknown mode", func() {
			_, err := getTplCfg("test", false)
			Ω(err).Should(MatchError("command is not supported"))
		})
	})
})
