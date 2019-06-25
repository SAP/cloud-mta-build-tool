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
		initCmdSrc = getTestPath("mta")
		initCmdTrg = getTestPath("result")
		initCmd.Run(nil, []string{})
		Ω(getTestPath("result", "Makefile.mta")).Should(BeAnExistingFile())
	})
})

var _ = Describe("Build", func() {

	BeforeEach(func() {
		os.Mkdir(getTestPath("result"), os.ModePerm)
	})
	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	It("Failure - wrong platform", func() {
		buildProjectCmdSrc = getTestPath("mta")
		buildProjectCmdTrg = getTestPath("result")
		buildProjectCmdPlatform = "xxx"
		err := buildCmd.RunE(nil, []string{})
		Ω(err).ShouldNot(BeNil())
	})
})
