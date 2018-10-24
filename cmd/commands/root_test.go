package commands

import (
	"cloud-mta-build-tool/internal/fsys"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
	"github.com/spf13/viper"
)

var _ = Describe("Root", func() {

	Describe("Init config", func() {
		AfterEach(func() {
			viper.Reset()
			cfgFile = ""
		})

		It("config file not defined", func() {
			viper.Reset()
			Ω(viper.Get("xxx")).Should(BeNil())
		})

		DescribeTable("config file defined", func(configFilename string, matcher GomegaMatcher) {
			cfgFile, _ = dir.GetFullPath("testdata", configFilename)
			initConfig()
			Ω(viper.Get("xxx")).Should(matcher)
		},
			Entry("right config", "config.props", Equal("10")),
			Entry("wrong config", "config1.props", BeNil()),
		)
	})

	Describe("Execute", func() {
		It("Sanity", func() {
			out := executeAndProvideOutput(func() {
				Execute()
			})
			Ω(out).Should(ContainSubstring("help"))
		})
	})

})
