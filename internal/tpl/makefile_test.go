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
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
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
		tpl              = tplCfg{tplContent: makeVerbose, relPath: "", preContent: basePreVerbose, postContent: basePostVerbose, depDesc: "dev"}
		tplDep           = tplCfg{tplContent: makeVerboseDep, relPath: "", preContent: basePreVerbose, postContent: basePostVerbose, depDesc: "dep"}
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

		var _ = Describe("ExecuteMake", func() {
			AfterEach(func() {
				os.Remove(filepath.Join(wd, "testdata", "Makefile.mta"))
			})
			It("Sanity", func() {
				Ω(ExecuteMake(filepath.Join(wd, "testdata"), filepath.Join(wd, "testdata"), "dev", "", os.Getwd)).Should(Succeed())
				Ω(filepath.Join(wd, "testdata", "Makefile.mta")).Should(BeAnExistingFile())
			})
			It("Fails on location initialization", func() {
				Ω(ExecuteMake("", filepath.Join(wd, "testdata"), "dev", "", func() (string, error) {
					return "", errors.New("err")
				})).Should(HaveOccurred())
			})
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
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), Descriptor: "dev"}
			Ω(makeFile(&ep, &ep, makeFileName, &tpl)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileContent))
		})
		It("Sanity - Dep", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), Descriptor: "dep"}
			Ω(makeFile(&ep, &ep, makeFileName, &tplDep)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileDepContent))
		})
		It("genMakefile testing with wrong mta yaml file", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), MtaFilename: "xxx.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, "")).Should(HaveOccurred())
		})
		It("genMakefile testing with wrong mode", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			Ω(genMakefile(&ep, &ep, &ep, "wrongMode")).Should(HaveOccurred())
		})
	})

	var _ = DescribeTable("Makefile Generation Failed", func(testPath string, testTemplateFilename string) {
		wd, _ := os.Getwd()
		testTemplate, _ := ioutil.ReadFile(filepath.Join(wd, "testdata", testTemplateFilename))
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata")}
		Ω(makeFile(&ep, &ep, makeFileName, &tplCfg{relPath: testPath, tplContent: testTemplate, preContent: basePreVerbose, postContent: basePostVerbose, depDesc: "dev"})).Should(HaveOccurred())
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

	var _ = Describe("genMakefile mode tests", func() {
		DescribeTable("Positive", func(mode string, tpl tplCfg, isDep bool) {
			Ω(getTplCfg(mode, isDep)).Should(Equal(tpl))
		},
			Entry("Default mode Dev", "", tplCfg{tplContent: makeDefault, preContent: basePreDefault, postContent: basePostDefault}, false),
			Entry("Default mode Dep", "", tplCfg{tplContent: makeDefault, preContent: basePreDefault, postContent: basePostDefault}, true),
			Entry("Verbose mode Dev", "verbose", tplCfg{tplContent: makeVerbose, preContent: basePreVerbose, postContent: basePostVerbose}, false),
			Entry("Verbose mode Dep", "verbose", tplCfg{tplContent: makeVerboseDep, preContent: basePreVerbose, postContent: basePostVerbose}, true),
		)
		It("unknown mode", func() {
			_, err := getTplCfg("test", false)
			Ω(err).Should(MatchError("command is not supported"))
		})
	})
})
