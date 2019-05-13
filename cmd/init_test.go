package commands

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
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

var _ = Describe("Build", func() {

	BeforeEach(func() {
		os.Mkdir(getTestPath("result"), os.ModePerm)
	})
	AfterEach(func() {
		//os.RemoveAll(getTestPath("result"))
	})
	It("Sanity", func() {
		buildProjectCmdDesc = "dev"

		buildProjectCmdSrc = getTestPath("mta")
		buildProjectCmdTrg = getTestPath("result")
		buildProjectCmdPlatform = "cf"

		buildCmd.Run(nil, []string{})
		//Ω(err).Should(Succeed())
	})
	It("Invalid descriptor", func() {
		initCmdDesc = "xx"
		initCmdSrc = getTestPath("mta")
		initCmdTrg = getTestPath("result")
		initCmd.Run(nil, []string{})
		Ω(getTestPath("result", "Makefile_tmp.mta")).ShouldNot(BeAnExistingFile())
	})
})
