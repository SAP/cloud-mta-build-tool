package commands

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

var _ = Describe("DeprecatedBuilders", func() {

	var originalDeprecatedBuilders = deprecatedBuilders

	BeforeEach(func() {
		deprecatedBuilders = map[string]string{"deprecated_builder": `the "deprecated_builder" builder is deprecated and will be removed soon.`}
		logs.Logger = logs.NewLogger()
	})

	AfterEach(func() {
		deprecatedBuilders = originalDeprecatedBuilders
	})

	var _ = Describe("awareOfDeprecatedBuilder function", func() {
		It("Does not log warning in case of not deprecated builder", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for warnings analysis
			logs.Logger.SetOutput(&str)
			awareOfDeprecatedBuilder("new_builder")
			Ω(str.String()).Should(BeEmpty())
		})

		It("Logs warning in case of deprecated builder", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for warnings analysis
			logs.Logger.SetOutput(&str)
			awareOfDeprecatedBuilder("deprecated_builder")
			Ω(str.String()).Should(ContainSubstring(deprecatedBuilders["deprecated_builder"]))
		})

	})
})
