package artifacts

import (
	"fmt"
	"os"

	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta/mta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {

	var _ = Describe("ExecuteProjectBuild", func() {
		It("wrong phase", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "dev", "wrong phase", os.Getwd)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal(`the "wrong phase" phase of mta project build is invalid; supported phases: "pre", "post"`))
		})
		It("wrong location", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "xx", "pre", func() (string, error) {
				return "", fmt.Errorf("error")
			})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("failed to initialize the location when validating descriptor:"))
		})
		It("mta.yaml not found", func() {
			err := ExecuteProjectBuild(getTestPath("mta1"), "dev", "pre", os.Getwd)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("failed to read"))
		})
		It("Sanity - wrong builder", func() {
			err := ExecuteProjectBuild(getTestPath("mta"), "dev", "pre", os.Getwd)
			Ω(err).Should(HaveOccurred())
		})
	})

	var _ = Describe("getProjectBuilderCommands", func() {
		It("Builder and commands defined", func() {
			projectBuild := mta.ProjectBuilder{
				Builder: "npm",
				Options: mta.ProjectBuilderOptions{
					Execute: []string{"command {{xxx.abc}}"},
				},
				BuildParams: map[string]interface{}{
					"xxx-opts": map[interface{}]interface{}{
						"abc": "aaa",
					},
					"npm-opts": map[interface{}]interface{}{
						"config": map[interface{}]interface{}{
							"foo": "xyz",
						},
					},
				},
			}

			cmds, err := getProjectBuilderCommands(projectBuild)
			Ω(err).Should(Succeed())
			Ω(len(cmds.Command)).Should(Equal(3))
			Ω(cmds.Command[0]).Should(Equal("npm install  --foo xyz"))
			Ω(cmds.Command[2]).Should(Equal("command aaa"))
		})
	})

	var _ = Describe("execProjectBuilders", func() {
		It("Before Defined with nothing to execute", func() {
			builders := []mta.ProjectBuilder{}
			projectBuild := mta.ProjectBuild{
				BeforeAll: struct {
					Builders []mta.ProjectBuilder `yaml:"builders,omitempty"`
				}{Builders: builders},
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}
			Ω(execProjectBuilders(&oMta, "pre")).Should(Succeed())
		})
		It("After Defined with nothing to execute", func() {
			builders := []mta.ProjectBuilder{}
			projectBuild := mta.ProjectBuild{
				AfterAll: struct {
					Builders []mta.ProjectBuilder `yaml:"builders,omitempty"`
				}{Builders: builders},
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}
			Ω(execProjectBuilders(&oMta, "post")).Should(Succeed())
		})
		It("Before Defined with wrong builder", func() {
			builders := []mta.ProjectBuilder{
				{
					Builder: "xxx",
				},
			}
			projectBuild := mta.ProjectBuild{
				BeforeAll: struct {
					Builders []mta.ProjectBuilder `yaml:"builders,omitempty"`
				}{Builders: builders},
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}
			Ω(execProjectBuilders(&oMta, "pre")).Should(HaveOccurred())
		})
		It("After Defined with wrong builder", func() {
			builders := []mta.ProjectBuilder{
				{
					Builder: "xxx",
				},
			}
			projectBuild := mta.ProjectBuild{
				AfterAll: struct {
					Builders []mta.ProjectBuilder `yaml:"builders,omitempty"`
				}{Builders: builders},
			}
			oMta := mta.MTA{
				BuildParams: &projectBuild,
			}
			Ω(execProjectBuilders(&oMta, "post")).Should(HaveOccurred())
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
			Ω(execBuilder([]mta.ProjectBuilder{builder})).Should(Succeed())
			commands.BuilderTypeConfig = buildersCfg
		})
		It("Builder does not exist", func() {
			builder := mta.ProjectBuilder{
				Builder: "testbuilder",
			}
			Ω(execBuilder([]mta.ProjectBuilder{builder})).Should(HaveOccurred())
		})
		It("Sanity - no builder defined", func() {

			builder := mta.ProjectBuilder{}
			Ω(execBuilder([]mta.ProjectBuilder{builder})).Should(Succeed())
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
			Ω(execBuilder([]mta.ProjectBuilder{builder})).Should(HaveOccurred())
			commands.BuilderTypeConfig = buildersCfg
		})
	})

})
