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
		initCmdDesc = "dev"
		initCmdSrc = getTestPath("mta")
		initCmdTrg = getTestPath("result")
		initCmd.Run(nil, []string{})
		Ω(getTestPath("result", "Makefile.mta")).Should(BeAnExistingFile())
	})
	It("Invalid descriptor", func() {
		initCmdDesc = "xx"
		initCmdSrc = getTestPath("mta")
		initCmdTrg = getTestPath("result")
		initCmd.Run(nil, []string{})
		Ω(getTestPath("result", "Makefile.mta")).ShouldNot(BeAnExistingFile())
	})
})
