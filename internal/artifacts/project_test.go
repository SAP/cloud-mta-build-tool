package artifacts

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Project", func() {

	var _ = Describe("ExecuteProjectBuild", func() {
		It("Sanity - post phase", func() {
			err := ExecuteProjectBuild(getTestPath("mtahtml5"), "dev", "post", os.Getwd)
			Ω(err).Should(BeNil())
		})
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
			Ω(execBuilder("testbuilder")).Should(Succeed())
			commands.BuilderTypeConfig = buildersCfg
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
			Ω(execBuilder("testbuilder")).Should(HaveOccurred())
			commands.BuilderTypeConfig = buildersCfg
		})
		Context("pre & post builder commands", func() {
			oMta := &mta.MTA{}
			BeforeEach(func() {
				mtaFile, _ := ioutil.ReadFile("./testdata/mta/mta.yaml")
				yaml.Unmarshal(mtaFile, oMta)
			})
			It("before-all builder", func() {
				v := beforeExec(oMta.BuildParams)
				Ω(v).Should(Equal("mybuilder"))
			})

			It("after-all builder", func() {
				v := afterExec(oMta.BuildParams)
				Ω(v).Should(Equal("otherbuilder"))
			})
		})
		Context("pre & post builder commands - no builders defined", func() {
			oMta := &mta.MTA{}
			BeforeEach(func() {
				mtaFile, _ := ioutil.ReadFile("./testdata/mta/mta.yaml")
				yaml.Unmarshal(mtaFile, oMta)
				oMta.BuildParams.BeforeAll.Builders = nil
				oMta.BuildParams.AfterAll.Builders = nil
			})
			It("before-all builder", func() {
				v := beforeExec(oMta.BuildParams)
				Ω(v).Should(Equal(""))
			})

			It("after-all builder", func() {
				v := afterExec(oMta.BuildParams)
				Ω(v).Should(Equal(""))
			})
		})
	})

})
