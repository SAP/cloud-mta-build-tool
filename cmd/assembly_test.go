package commands

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

var _ = Describe("Assembly", func() {

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	It("Sanity", func() {
		err := artifacts.Assembly(getTestPath("assembly-sample"),
			getTestPath("result"), "cf", "", "?", os.Getwd)
		Ω(err).Should(Succeed())
		Ω(getTestPath("result", "com.sap.xs2.samples.javahelloworld_0.1.0.mtar")).Should(BeAnExistingFile())
	})
	var _ = DescribeTable("Fails on location initialization", func(maxCalls int) {
		calls := 0
		err := artifacts.Assembly("",
			getTestPath("result"), "cf", "", "true", func() (string, error) {
				calls++
				if calls >= maxCalls {
					return "", errors.New("error")
				}
				return getTestPath("assembly-sample"), nil
			})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("failed to initialize the location when getting working directory"))
	},
		Entry("fails on CopyMtaContent", 1),
		Entry("fails on ExecuteGenMeta", 2),
		Entry("fails on ExecuteGenMtar", 3),
		Entry("fails on ExecuteCleanup", 4),
	)

	It("assemblyCommand - fails on missing mtad in the current location", func() {
		Ω(assemblyCommand.RunE(nil, []string{})).Should(HaveOccurred())
	})

})
