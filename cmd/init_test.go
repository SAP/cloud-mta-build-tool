package commands

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {

	BeforeEach(func() {
		err := os.Mkdir(getTestPath("result"), os.ModePerm)
		if err != nil {
			fmt.Println("error occurred during dir creation")
		}
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
		err := os.Mkdir(getTestPath("result"), os.ModePerm)
		if err != nil {
			fmt.Println("error occurred during dir creation")
		}
	})
	AfterEach(func() {
		err := os.RemoveAll(getTestPath("result"))
		if err != nil {
			fmt.Println("error occurred during dir cleanup")
		}
	})
	It("Failure - wrong platform", func() {
		buildCmdSrc = getTestPath("mta")
		buildCmdTrg = getTestPath("result")
		buildCmdPlatform = "xxx"
		err := buildCmd.RunE(nil, []string{})
		Ω(err).Should(HaveOccurred())
	})
})
