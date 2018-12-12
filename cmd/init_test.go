package commands

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {

	BeforeEach(func() {
		os.Mkdir(getTestPath("result"), os.ModePerm)
	})
	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	It("Sanity", func() {
		descriptorInitFlag = "dev"
		sourceInitFlag = getTestPath("mta")
		targetInitFlag = getTestPath("result")
		initProcessCmd.Run(nil, []string{})
		Ω(getTestPath("result", "Makefile.mta")).Should(BeAnExistingFile())
	})
	It("Invalid descriptor", func() {
		descriptorInitFlag = "xx"
		sourceInitFlag = getTestPath("mta")
		targetInitFlag = getTestPath("result")
		initProcessCmd.Run(nil, []string{})
		Ω(getTestPath("result", "Makefile.mta")).ShouldNot(BeAnExistingFile())
	})
})
