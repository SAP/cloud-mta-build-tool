package commands

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {

	BeforeEach(func() {
		Ω(os.MkdirAll(getTestPath("result"), os.ModePerm)).Should(Succeed())
	})
	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
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
		Ω(os.MkdirAll(getTestPath("result"), os.ModePerm)).Should(Succeed())
	})
	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})
	It("Failure - wrong platform", func() {
		buildCmdSrc = getTestPath("mta")
		buildCmdTrg = getTestPath("result")
		buildCmdPlatform = "xxx"
		err := buildCmd.RunE(nil, []string{})
		Ω(err).Should(HaveOccurred())
	})
})
