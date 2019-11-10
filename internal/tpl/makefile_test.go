package tpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kballard/go-shellquote"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta/mta"
)

const (
	makefile = "Makefile.mta"
)

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})

func removeSpecialSymbols(b []byte) string {
	s := string(b)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func getMakeFileContent(filePath string) string {
	expected, _ := ioutil.ReadFile(filePath)
	return removeSpecialSymbols(expected)
}

func escapePath(parts ...string) string {
	return shellquote.Join(filepath.Join(parts...))
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
cli_version: 0.0.0
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
			Ω(makeFile(&ep, &ep, &ep, nil, makeFileName, &tpl, true)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			Ω(getMakeFileContent(makeFileFullPath)).Should(Equal(expectedMakeFileContent))
		})
		It("Create make file in folder that does not exist", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata", "someFolder"), Descriptor: "dev"}
			Ω(makeFile(&ep, &ep, &ep, nil, makeFileName, &tpl, true)).Should(Succeed())
			filename := filepath.Join(ep.GetTarget(), makeFileName)
			Ω(filename).Should(BeAnExistingFile())
			Ω(getMakeFileContent(filename)).Should(Equal(expectedMakeFileContent))
		})
		It("genMakefile testing with wrong mta yaml file", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), MtaFilename: "xxx.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, &ep, nil, makefile, "", true)).Should(HaveOccurred())
		})
		It("genMakefile testing with wrong target folder (file path)", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata", "mta.yaml"), MtaFilename: "xxx.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, &ep, nil, makefile, "", true)).Should(HaveOccurred())
		})
		It("genMakefile testing with wrong mode", func() {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata")}
			Ω(genMakefile(&ep, &ep, &ep, &ep, nil, makefile, "wrongMode", true)).Should(HaveOccurred())
		})

		DescribeTable("genMakefile should fail when there is a circular build dependency between modules", func(mode string) {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata"), MtaFilename: "circular.yaml"}
			Ω(genMakefile(&ep, &ep, &ep, &ep, nil, makefile, mode, true)).Should(HaveOccurred())
		},
			Entry("in default mode", ""),
			Entry("in verbose mode", "verbose"),
		)

		DescribeTable("generate module build in verbose make file", func(mtaFileName, moduleName, expectedModuleCommandsGen string) {
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata", "modulegen"), TargetPath: filepath.Join(wd, "testdata"), Descriptor: "dev", MtaFilename: mtaFileName}
			Ω(makeFile(&ep, &ep, &ep, nil, makeFileName, &tpl, true)).Should(Succeed())
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

		modulegen := filepath.Join(wd, "testdata", "modulegen")
		DescribeTable("generate module build with dependencies in verbose make file", func(mtaFileName, moduleName, modulePath, expectedModuleDepNames string, expectedModuleDepCopyCommands string) {
			ep := dir.Loc{SourcePath: modulegen, TargetPath: filepath.Join(wd, "testdata"), Descriptor: "dev", MtaFilename: mtaFileName}
			Ω(makeFile(&ep, &ep, &ep, nil, makeFileName, &tpl, true)).Should(Succeed())
			Ω(makeFileFullPath).Should(BeAnExistingFile())
			makefileContent := getMakeFileContent(makeFileFullPath)

			expectedModuleGen := fmt.Sprintf(`%s: validate %s%s
	@cd "$(PROJ_DIR)/%s" &&`, moduleName, expectedModuleDepNames, expectedModuleDepCopyCommands, modulePath)
			Ω(makefileContent).Should(ContainSubstring(removeSpecialSymbols([]byte(expectedModuleGen))))
		},
			Entry("dependency with artifacts", "dep_with_patterns.yaml", "module1", "public", `dep`, fmt.Sprintf(`
	@$(MBT) cp -s=%s -t=%s -p=dist/\* -p=some_dir -p=a\*.txt`, escapePath(modulegen, "client"), escapePath(modulegen, "public"))),
			Entry("module with two dependencies", "two_deps.yaml", "my_proj_ui_deployer", "my_proj_ui_deployer", `ui5module1 ui5module2`, fmt.Sprintf(`
	@$(MBT) cp -s=%s -t=%s -p=./\*
	@$(MBT) cp -s=%s -t=%s -p=./\*`,
				escapePath(modulegen, "ui5module1", "dist"), escapePath(modulegen, "my_proj_ui_deployer", "resources", "ui5module1"),
				escapePath(modulegen, "ui5module2", "dist"), escapePath(modulegen, "my_proj_ui_deployer", "resources", "ui5module2"))),
			Entry("dependency with target-path", "dep_with_artifacts_and_targetpath.yaml", "module1", "public", `module1-dep`, fmt.Sprintf(`
	@$(MBT) cp -s=%s -t=%s -p=dist/\*`, escapePath(modulegen, "client"), escapePath(modulegen, "public", "client"))),
			Entry("dependent module with build-result and module with artifacts and target-path", "dep_with_build_results.yaml", "module1", "public", `dep1 dep2`, fmt.Sprintf(`
	@$(MBT) cp -s=%s -t=%s -p=\*
	@$(MBT) cp -s=%s -t=%s -p=\*`,
				escapePath(modulegen, "client1", "dist"), escapePath(modulegen, "public", "dep1_result"),
				escapePath(modulegen, "client2", "target/*.war"), escapePath(modulegen, "public"))),
		)
	})

	DescribeTable("Makefile Generation Failed", func(testPath string, testTemplateFilename string) {
		wd, _ := os.Getwd()
		testTemplate, _ := ioutil.ReadFile(filepath.Join(wd, "testdata", testTemplateFilename))
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), TargetPath: filepath.Join(wd, "testdata")}
		Ω(makeFile(&ep, &ep, &ep, nil, makeFileName, &tplCfg{relPath: testPath, tplContent: testTemplate, preContent: basePreVerbose, postContent: basePost, depDesc: "dev"}, true)).Should(HaveOccurred())
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

	var absPath = func(path string) string {
		s, _ := filepath.Abs(path)
		return s
	}
	sep := string(filepath.Separator)
	DescribeTable("getExtensionsArg", func(extensions []string, makefileDirPath string, expected string) {
		Ω(getExtensionsArg(extensions, makefileDirPath, "-e")).Should(Equal(expected))
	},
		Entry("empty list returns empty string", []string{}, "", ""),
		Entry("nil returns empty string", nil, "", ""),
		Entry("extension path is returned relative to the makefile path when it's in the same folder",
			[]string{absPath("my.mtaext")}, absPath("."), ` -e="$(CURDIR)`+sep+`my.mtaext"`),
		Entry("extension path is returned relative to the makefile path when it's in an inner folder",
			[]string{absPath(filepath.Join("inner", "my.mtaext"))}, absPath("."), ` -e="$(CURDIR)`+sep+"inner"+sep+`my.mtaext"`),
		Entry("extension path is returned relative to the makefile path when it's in an outer folder",
			[]string{absPath("my.mtaext")}, absPath("inner"), ` -e="$(CURDIR)`+sep+".."+sep+`my.mtaext"`),
		Entry("extension paths are separated by a comma",
			[]string{absPath("my.mtaext"), absPath("second.mtaext")}, absPath("."), ` -e="$(CURDIR)`+sep+`my.mtaext,$(CURDIR)`+sep+`second.mtaext"`),
	)

	It("GetModuleDeps returns error when module doesn't exist", func() {
		data := templateData{File: mta.MTA{}}
		_, err := data.GetModuleDeps("unknown")
		Ω(err).Should(HaveOccurred())
	})

	It("GetModuleDeps returns error when module has dependency that doesn't exist", func() {
		data := templateData{File: mta.MTA{Modules: []*mta.Module{
			{
				Name: "m1",
				BuildParams: map[string]interface{}{
					"requires": []interface{}{
						map[string]interface{}{"name": "unknown"},
					},
				},
			},
		}}}
		_, err := data.GetModuleDeps("m1")
		Ω(err).Should(HaveOccurred())
	})
})
