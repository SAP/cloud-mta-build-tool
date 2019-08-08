package tpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	s := string(b)
	fmt.Println(s)
	s = strings.Replace(s, "\r", "", -1)
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

	Describe("MakeFile Generation", func() {
		BeforeEach(func() {
			version.VersionConfig = []byte(`
cli_version: v0.0.0
makefile_version: 0.0.0
`)
		})
		AfterEach(func() {
			Ω(os.RemoveAll(makeFileFullPath)).Should(Succeed())
			Ω(os.RemoveAll(filepath.Join(wd, "testdata", "someFolder"))).Should(Succeed())
		})

		Describe("ExecuteMake", func() {
			AfterEach(func() {
				Ω(os.RemoveAll(filepath.Join(wd, "testdata", "Makefile.mta"))).Should(Succeed())
			})
			It("Sanity", func() {
				Ω(ExecuteMake(filepath.Join(wd, "testdata"), filepath.Join(wd, "testdata"), nil, makefile, "", os.Getwd, true)).Should(Succeed())
				Ω(filepath.Join(wd, "testdata", "Makefile.mta")).Should(BeAnExistingFile())
			})
			It("Fails on location initialization", func() {
				Ω(ExecuteMake("", filepath.Join(wd, "testdata"), nil, makefile, "", func() (string, error) {
					return "", errors.New("err")
				}, true)).Should(HaveOccurred())
			})
			It("Fails on wrong mode", func() {
				Ω(ExecuteMake(filepath.Join(wd, "testdata"), filepath.Join(wd, "testdata"), nil, makefile, "wrong", os.Getwd, true)).Should(HaveOccurred())
			})
		})

		It("createMakeFile testing", func() {
			makeFilePath := filepath.Join(wd, "testdata")
			file, _ := createMakeFile(makeFilePath, makeFileName)
			Ω(file).ShouldNot(BeNil())
			Ω(file.Close()).Should(Succeed())
			Ω(makeFilePath).Should(BeAnExistingFile())
			_, err := createMakeFile(makeFilePath, makeFileName)
			Ω(err).Should(HaveOccurred())
		})
		It("Sanity", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), Descriptor: "dev"}
			Ω(makeFile(&ep, &ep, nil, makeFileName, &tpl, true)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileContent))
		})
		It("Create make file in folder that does not exist", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata", "someFolder"), Descriptor: "dev"}
			Ω(makeFile(&ep, &ep, nil, makeFileName, &tpl, true)).Should(Succeed())
			filename := filepath.Join(ep.GetTarget(), makeFileName)
			Ω(filename).Should(BeAnExistingFile())
			Ω(getMakeFileContent(filename)).Should(Equal(expectedMakeFileContent))
		})
		It("genMakefile testing with wrong mta yaml file", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), MtaFilename: "xxx.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, nil, makefile, "", true)).Should(HaveOccurred())
		})
		It("genMakefile testing with wrong target folder (file path)", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata", "mta.yaml"), MtaFilename: "xxx.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, nil, makefile, "", true)).Should(HaveOccurred())
		})
		It("genMakefile testing with wrong mode", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			Ω(genMakefile(&ep, &ep, &ep, nil, makefile, "wrongMode", true)).Should(HaveOccurred())
		})

		DescribeTable("generate module build in verbose make file", func(mtaFileName, moduleName, expectedModuleCommandsGen string) {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata", "modulegen"), TargetPath: filepath.Join(wd, "testdata"), Descriptor: "dev", MtaFilename: mtaFileName}
			Ω(makeFile(&ep, &ep, nil, makeFileName, &tpl, true)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			makefileContent := getMakeFileContent(makeFileFullPath)

			expectedModuleGen := fmt.Sprintf(`%s: validate
	@cd "$(PROJ_DIR)/%s" && %s`, moduleName, moduleName, expectedModuleCommandsGen)
			Ω(makefileContent).Should(ContainSubstring(removeSpecialSymbols([]byte(expectedModuleGen))))
		},
			Entry("module with one command", "one_command.yaml", "one_command", `$(MBT) execute -c=yarn`),
			Entry("module with no commands and no timeout",
				"no_commands.yaml", "no_commands", `$(MBT) execute`),
			Entry("module with no commands and with timeout",
				"no_commands_with_timeout.yaml", "no_commands_with_timeout", `$(MBT) execute -t=3m`),
			Entry("module with multiple commands",
				"multiple_commands.yaml", "multiple_commands", `$(MBT) execute -c='npm install' -c=grunt -c='npm prune --production'`),
			Entry("module with command and timeout",
				"command_with_timeout.yaml", "command_with_timeout", `$(MBT) execute -t=2s -c='sleep 1'`),
			Entry("module with commands with special characters",
				"commands_with_special_chars.yaml", "commands_with_special_chars", `$(MBT) execute -c='bash -c '\''echo "a"'\' -c='echo "a\b"'`),
		)
	})

	DescribeTable("Makefile Generation Failed", func(testPath string, testTemplateFilename string) {
		wd, _ := os.Getwd()
		testTemplate, _ := ioutil.ReadFile(filepath.Join(wd, "testdata", testTemplateFilename))
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata")}
		Ω(makeFile(&ep, &ep, nil, makeFileName, &tplCfg{relPath: testPath, tplContent: testTemplate, preContent: basePreVerbose, postContent: basePost, depDesc: "dev"}, true)).Should(HaveOccurred())
	},
		Entry("Wrong Template", "testdata", filepath.Join("testdata", "WrongMakeTmpl.txt")),
		Entry("Yaml not exists", "testdata1", "make_default.txt"),
	)

	DescribeTable("String in slice search", func(s string, slice []string, expected bool) {
		Ω(stringInSlice(s, slice)).Should(Equal(expected))
	},
		Entry("positive test", "test1", []string{"test1", "foo"}, true),
		Entry("negative test", "test1", []string{"--test", "foo"}, false),
	)

	Describe("genMakefile mode tests", func() {
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
