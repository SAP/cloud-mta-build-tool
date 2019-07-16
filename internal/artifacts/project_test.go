package artifacts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Project", func() {

	var _ = Describe("ExecuteProjectBuild", func() {
		It("Sanity - post phase", func() {
			err := ExecuteProjectBuild(getTestPath("mtahtml5"), "dev", "post", os.Getwd)
			Ω(err).Should(Succeed())
		})
		It("wrong phase", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "dev", "wrong phase", os.Getwd)
			checkError(err, UnsupportedPhaseMsg, "wrong phase")
		})
		It("wrong location", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "xx", "pre", func() (string, error) {
				return "", fmt.Errorf("error")
			})
			checkError(err, dir.InvalidDescMsg, "xx")
		})
		It("mta.yaml not found", func() {
			err := ExecuteProjectBuild(getTestPath("mta1"), "dev", "pre", os.Getwd)
			checkError(err, dir.ReadFailedMsg, getTestPath("mta1", "mta.yaml"))
		})
		It("Sanity - custom builder", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "dev", "pre", os.Getwd)
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
		})
		It("Sanity", func() {
			err := ExecBuild("Makefile_tmp.mta", getTestPath("mta_with_zipped_module"), getResultPath(), "", "", "cf", true, os.Getwd, func(strings [][]string) error {
				return nil
			})
			Ω(err).Should(Succeed())
			Ω(filepath.Join(getResultPath(), "Makefile_tmp.mta")).ShouldNot(BeAnExistingFile())
		})
		It("Wrong - no platform", func() {
			err := ExecBuild("Makefile_tmp.mta", getTestPath("mta_with_zipped_module"), getResultPath(), "", "", "", true, os.Getwd, func(strings [][]string) error {
				return fmt.Errorf("failure")
			})
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
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`bash -c 'echo "aaa"'`}}
			Ω(execProjectBuilder([]mta.ProjectBuilder{builder}, "pre")).Should(Succeed())
		})

		It("Succeeds on builder with timeout, when timeout isn't reached", func() {
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`bash -c 'sleep 1'`}, Timeout: "10s"}
			Ω(execProjectBuilder([]mta.ProjectBuilder{builder}, "pre")).Should(Succeed())
		})

		It("Fails on builder with timeout, when timeout is reached", func() {
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`bash -c 'sleep 10'`}, Timeout: "2s"}
			err := execProjectBuilder([]mta.ProjectBuilder{builder}, "post")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(exec.ExecTimeoutMsg, "2s")))
		})

		It("Fails on builder with invalid custom command", func() {
			builder := mta.ProjectBuilder{Builder: "custom", Commands: []string{`bash -c 'sleep 10`}}
			err := execProjectBuilder([]mta.ProjectBuilder{builder}, "post")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(commands.BadCommandMsg, `bash -c 'sleep 10`)))
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
			oMta := &mta.MTA{}
			BeforeEach(func() {
				mtaFile, _ := ioutil.ReadFile("./testdata/mta/mta.yaml")
				Ω(yaml.Unmarshal(mtaFile, oMta)).Should(Succeed())
			})
		})
		Context("pre & post builder commands - no builders defined", func() {
			oMta := &mta.MTA{}
			BeforeEach(func() {
				mtaFile, _ := ioutil.ReadFile("./testdata/mta/mta.yaml")
				Ω(yaml.Unmarshal(mtaFile, oMta)).Should(Succeed())
				oMta.BuildParams.BeforeAll = nil
				oMta.BuildParams.AfterAll = nil
			})
		})
	})

})
