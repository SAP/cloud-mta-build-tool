package tpl

import (
	"fmt"
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

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
)

const (
	makefile = "Makefile.mta"
)

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})

func removeSpecialSymbols(b []byte) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	s := string(b)
	fmt.Println(s)
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
		tpl              = tplCfg{tplContent: makeVerbose, relPath: "", preContent: basePreVerbose, postContent: basePost, depDesc: "dev"}
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
		makeFileFullPath = func() string {
			return filepath.Join(wd, "testdata", makeFileName)
		}()
		expectedMakeFileContent = getMakeFileContent(expectedMakePath)
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
			os.RemoveAll(filepath.Join(wd, "testdata", "someFolder"))
		})

		var _ = Describe("ExecuteMake", func() {
			AfterEach(func() {
				os.Remove(filepath.Join(wd, "testdata", "Makefile.mta"))
			})
			It("Sanity", func() {
				Ω(ExecuteMake(filepath.Join(wd, "testdata"), filepath.Join(wd, "testdata"), makefile, "", os.Getwd)).Should(Succeed())
				Ω(filepath.Join(wd, "testdata", "Makefile.mta")).Should(BeAnExistingFile())
			})
			It("Fails on location initialization", func() {
				Ω(ExecuteMake("", filepath.Join(wd, "testdata"), makefile, "", func() (string, error) {
					return "", errors.New("err")
				})).Should(HaveOccurred())
			})
			It("Fails on wrong mode", func() {
				Ω(ExecuteMake(filepath.Join(wd, "testdata"), filepath.Join(wd, "testdata"), makefile, "wrong", os.Getwd)).Should(HaveOccurred())
			})
		})

		It("createMakeFile testing", func() {
			makeFilePath := filepath.Join(wd, "testdata")
			file, _ := createMakeFile(makeFilePath, makeFileName)
			Ω(file).ShouldNot(BeNil())
			file.Close()
			Ω(makeFilePath).Should(BeAnExistingFile())
			_, err := createMakeFile(makeFilePath, makeFileName)
			Ω(err).Should(HaveOccurred())
		})
		It("Sanity", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), Descriptor: "dev"}
			Ω(makeFile(&ep, &ep, makeFileName, &tpl)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileContent))
		})
		It("Create make file in folder that does not exist", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata", "someFolder"), Descriptor: "dev"}
			Ω(makeFile(&ep, &ep, makeFileName, &tpl)).Should(Succeed())
			filename := filepath.Join(ep.GetTarget(), makeFileName)
			Ω(filename).Should(BeAnExistingFile())
			Ω(getMakeFileContent(filename)).Should(Equal(expectedMakeFileContent))
		})
		It("genMakefile testing with wrong mta yaml file", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), MtaFilename: "xxx.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, makefile, "")).Should(HaveOccurred())
		})
		It("genMakefile testing with wrong target folder (file path)", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata", "mta.yaml"), MtaFilename: "xxx.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, makefile, "")).Should(HaveOccurred())
		})
		It("genMakefile testing with wrong mode", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			Ω(genMakefile(&ep, &ep, &ep, makefile, "wrongMode")).Should(HaveOccurred())
		})
	})

	var _ = DescribeTable("Makefile Generation Failed", func(testPath string, testTemplateFilename string) {
		wd, _ := os.Getwd()
		testTemplate, _ := ioutil.ReadFile(filepath.Join(wd, "testdata", testTemplateFilename))
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata")}
		Ω(makeFile(&ep, &ep, makeFileName, &tplCfg{relPath: testPath, tplContent: testTemplate, preContent: basePreVerbose, postContent: basePost, depDesc: "dev"})).Should(HaveOccurred())
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
			Entry("Default mode Dev", "", tplCfg{tplContent: makeDefault, preContent: basePreDefault, postContent: basePost}, false),
			Entry("Verbose mode Dev", "verbose", tplCfg{tplContent: makeVerbose, preContent: basePreVerbose, postContent: basePost}, false),
		)
		It("unknown mode", func() {
			_, err := getTplCfg("test", false)
			Ω(err).Should(MatchError(`the "test" command is not supported`))
		})
	})
})
