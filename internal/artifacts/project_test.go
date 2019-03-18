package artifacts

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Project", func() {
	var _ = Describe("runBuilder", func() {
		It("Sanity", func() {
			buildersCfg := commands.BuilderTypeConfig
			commands.BuilderTypeConfig =

				[]byte(`
builders:
- name: testbuilder
  info: "installing module dependencies & remove dev dependencies"
  path: "path to config file which override the following default commands"
  type:
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
  type:
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
				v := beforeExec(oMta)
				Ω(v).Should(Equal("mybuilder"))
			})

			It("after-all builder", func() {
				v := afterExec(oMta)
				Ω(v).Should(Equal("otherbuilder"))
			})
		})
	})

})
