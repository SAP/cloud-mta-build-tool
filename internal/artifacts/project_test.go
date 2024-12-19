package artifacts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Project", func() {

	var _ = Describe("ExecuteProjectBuild", func() {
		It("Sanity - post phase", func() {
			err := ExecuteProjectBuild(getTestPath("mtahtml5"), "", "", "dev", nil, "post", os.Getwd)
			Ω(err).Should(Succeed())
		})
		It("wrong phase", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "", "", "dev", nil, "wrong phase", os.Getwd)
			checkError(err, UnsupportedPhaseMsg, "wrong phase")
		})
		It("wrong location", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "", "", "xx", nil, "pre", func() (string, error) {
				return "", fmt.Errorf("error")
			})
			checkError(err, dir.InvalidDescMsg, "xx")
		})
		It("mta.yaml not found", func() {
			err := ExecuteProjectBuild(getTestPath("mta1"), "", "", "dev", nil, "pre", os.Getwd)
			checkError(err, getTestPath("mta1", "mta.yaml"))
		})
		It("Sanity - custom builder", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "", "", "dev", nil, "pre", os.Getwd)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`"command1"`))
			Ω(err.Error()).Should(ContainSubstring("failed"))
		})
	})

	var _ = Describe("ExecBuild", func() {
		BeforeEach(func() {
			Ω(os.Mkdir(getTestPath("result"), os.ModePerm)).Should(Succeed())
		})
		AfterEach(func() {
			Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
			Ω(os.RemoveAll(filepath.Join(getTestPath("mta_with_zipped_module"), "Makefile_tmp.mta"))).Should(Succeed())
		})
		It("Sanity", func() {
			err := ExecBuild("Makefile_tmp.mta", getTestPath("mta_with_zipped_module"), "", getResultPath(), nil, "", "", "cf", true, 0, false, os.Getwd, func(strings [][]string, b bool) error {
				return nil
			}, true, false, "")
			Ω(err).Should(Succeed())
			Ω(filepath.Join(getTestPath("mta_with_zipped_module"), "Makefile_tmp.mta")).ShouldNot(BeAnExistingFile())
		})
		It("Sanity - keep makefile", func() {
			err := ExecBuild("Makefile_tmp.mta", getTestPath("mta_with_zipped_module"), "", getResultPath(), nil, "", "", "cf", true, 0, false, os.Getwd, func(strings [][]string, b bool) error {
				return nil
			}, true, true, "")
			Ω(err).Should(Succeed())
			Ω(filepath.Join(getTestPath("mta_with_zipped_module"), "Makefile_tmp.mta")).Should(BeAnExistingFile())
		})
		It("Wrong - no platform", func() {
			err := ExecBuild("Makefile_tmp.mta", getTestPath("mta_with_zipped_module"), "", getResultPath(), nil, "", "", "", true, 0, false, os.Getwd, func(strings [][]string, b bool) error {
				return fmt.Errorf("failure")
			}, true, false, "")
			Ω(err).Should(HaveOccurred())
		})
		It("Wrong - ExecuteMake fails on wrong location", func() {
			err := ExecBuild("Makefile_tmp.mta", "", "", getResultPath(), nil, "", "", "", true, 0, false,
				func() (string, error) {
					return "", errors.New("wrong location")
				}, func(strings [][]string, b bool) error {
					return nil
				}, true, false, "")
			Ω(err).Should(HaveOccurred())
		})
	})

	var _ = Describe("getProjectBuilderCommands", func() {
		It("Builder and commands defined", func() {
			projectBuild := mta.ProjectBuilder{
				Builder:  "npm",
				Commands: []string{"abc"},
			}

			cmds, err := getProjectBuilderCommands(projectBuild)
			Ω(err).Should(Succeed())
			Ω(len(cmds.Command)).Should(Equal(1))
			Ω(cmds.Command[0]).Should(Equal("npm install --production"))
		})
		It("Custom builder with no commands", func() {
			projectBuild := mta.ProjectBuilder{
				Builder: "custom",
			}

			cmds, err := getProjectBuilderCommands(projectBuild)
			Ω(err).Should(Succeed())
			Ω(len(cmds.Command)).Should(Equal(0))
		})
	})

	var _ = Describe("execProjectBuilders", func() {
		It("Before Defined with nothing to execute", func() {
			var builders []mta.ProjectBuilder
			projectBuild := mta.ProjectBuild{
				BeforeAll: builders,
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}

			Ω(execProjectBuilders(&dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}, &oMta, "pre")).Should(Succeed())
		})
		It("After Defined with nothing to execute", func() {
			var builders []mta.ProjectBuilder
			projectBuild := mta.ProjectBuild{
				AfterAll: builders,
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}
			Ω(execProjectBuilders(&dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}, &oMta, "post")).Should(Succeed())
		})
		It("Before Defined with wrong builder", func() {
			builders := []mta.ProjectBuilder{
				{
					Builder: "xxx",
				},
			}
			projectBuild := mta.ProjectBuild{
				BeforeAll: builders,
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}
			Ω(execProjectBuilders(&dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}, &oMta, "pre")).Should(HaveOccurred())
		})
		It("After Defined with wrong builder", func() {
			builders := []mta.ProjectBuilder{
				{
					Builder: "xxx",
				},
			}
			projectBuild := mta.ProjectBuild{
				AfterAll: builders,
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}
			Ω(execProjectBuilders(&dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}, &oMta, "post")).Should(HaveOccurred())
		})
	})

	var _ = Describe("runBuilder", func() {
		It("Sanity", func() {
			buildersCfg := commands.BuilderTypeConfig
			commands.BuilderTypeConfig =

				[]byte(`
builders:
- name: testbuilder
  info: "installing module dependencies & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands:
  - command: go version
`)
			builder := mta.ProjectBuilder{
				Builder: "testbuilder",
			}
			Ω(execProjectBuilder([]mta.ProjectBuilder{builder}, "pre")).Should(Succeed())
			commands.BuilderTypeConfig = buildersCfg
		})
		It("Builder does not exist", func() {
			builder := mta.ProjectBuilder{
				Builder: "testbuilder",
			}
			Ω(execProjectBuilder([]mta.ProjectBuilder{builder}, "pre")).Should(HaveOccurred())
		})
		It("Custom builder", func() {
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`sh -c 'echo "aaa"'`}}
			Ω(execProjectBuilder([]mta.ProjectBuilder{builder}, "pre")).Should(Succeed())
		})

		It("Succeeds on builder with timeout, when timeout isn't reached", func() {
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`sh -c 'sleep 1'`}, Timeout: "10s"}
			Ω(execProjectBuilder([]mta.ProjectBuilder{builder}, "pre")).Should(Succeed())
		})

		It("Fails on builder with timeout, when timeout is reached", func() {
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`sh -c 'sleep 10'`}, Timeout: "2s"}
			err := execProjectBuilder([]mta.ProjectBuilder{builder}, "post")
			checkError(err, exec.ExecTimeoutMsg, "2s")
		})

		It("Fails on builder with invalid custom command", func() {
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`sh -c 'sleep 10`}}
			err := execProjectBuilder([]mta.ProjectBuilder{builder}, "post")
			checkError(err, commands.BadCommandMsg, `sh -c 'sleep 10`)
		})

		It("Fails on command execution", func() {
			buildersCfg := commands.BuilderTypeConfig
			commands.BuilderTypeConfig =

				[]byte(`
builders:
- name: testbuilder
  info: "installing module dependencies & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands:
  - command: go test unknown_test.go
`)
			builder := mta.ProjectBuilder{
				Builder: "testbuilder",
			}
			Ω(execProjectBuilder([]mta.ProjectBuilder{builder}, "pre")).Should(HaveOccurred())
			commands.BuilderTypeConfig = buildersCfg
		})
		Context("pre & post builder commands", func() {
			It("parses pre and post commands", func() {
				mtaFile, _ := ioutil.ReadFile("./testdata/mta/mta.yaml")
				var err error
				_, err = mta.Unmarshal(mtaFile)
				Ω(err).Should(Succeed())
			})
		})
	})

	var _ = DescribeTable("createMakeCommand", func(target, mode string, strict bool, jobs int, cpus int, outputSync bool, additionalExpectedArgs []string) {
		command := createMakeCommand("Makefile_tmp", "./src", target, mode, "result.mtar", "cf", strict, jobs, outputSync, func() int {
			return cpus
		})
		Ω(len(command)).To(Equal(8+len(additionalExpectedArgs)), "number of command arguments")
		// The first arguments must be in this order
		Ω(command[0]).To(Equal("./src"))
		Ω(command[1]).To(Equal("make"))
		Ω(command[2]).To(Equal("-f"))
		Ω(command[3]).To(Equal("Makefile_tmp"))

		Ω(command).To(ContainElement("p=cf"))
		Ω(command).To(ContainElement("mtar=result.mtar"))
		Ω(command).To(ContainElement(fmt.Sprintf("strict=%v", strict)))
		Ω(command).To(ContainElement(fmt.Sprintf("mode=%s", mode)))

		for _, arg := range additionalExpectedArgs {
			Ω(command).To(ContainElement(arg))
		}
	},
		Entry("non-verbose without target", "", "", true, 0, 2, false, nil),
		Entry("non-verbose with target", "./trg", "", false, 0, 2, false, []string{`t="./trg"`}),
		Entry("non-verbose with specified jobs", "", "", true, 2, 4, false, nil),
		Entry("non-verbose with synchronized output", "", "", true, 2, 4, true, nil),
		Entry("verbose without target and without specified jobs, with less than the max number of CPUs", "", "verbose", false, 0, 2, false, []string{"-j2"}),
		Entry("verbose with target and without specified jobs, with more than the max number of CPUs", "./trg", "verbose", true, 0, MaxMakeParallel+10, false, []string{`t="./trg"`, fmt.Sprintf("-j%d", MaxMakeParallel)}),
		Entry("verbose with specified jobs less than the number of CPUs", "", "verbose", false, 3, 5, false, []string{"-j3"}),
		Entry("verbose with specified jobs less than the max number of CPUs", "", "verbose", false, 3, 20, false, []string{"-j3"}),
		Entry("verbose with specified jobs more than the number of CPUs", "", "v", true, 3, 1, false, []string{"-j3"}),
		Entry("verbose with specified jobs more than the max number of CPUs and less than the number of CPUs", "", "verbose", false, 20, 25, false, []string{"-j20"}),
		Entry("verbose with specified jobs more than the max number of CPUs and number of CPUs", "", "v", false, 20, 15, false, []string{"-j20"}),
		Entry("verbose with negative specified jobs", "", "verbose", true, -1, 3, false, []string{"-j3"}),
		Entry("verbose with negative specified jobs and more than the max number of CPUs", "", "verbose", true, -1, MaxMakeParallel+5, false, []string{fmt.Sprintf("-j%d", MaxMakeParallel)}),
		Entry("verbose without specified jobs and with synchronized output", "", "verbose", true, 0, 3, true, []string{"-j3", "-Otarget"}),
		Entry("verbose with one job and with synchronized output", "", "verbose", true, 1, 3, true, []string{"-j1", "-Otarget"}),
		Entry("verbose with several jobs and with synchronized output", "", "verbose", true, 2, 3, true, []string{"-j2", "-Otarget"}),
	)
})
