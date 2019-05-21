package commands

import (
	"os"
	"path/filepath"

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
	It("Failure - Makefile_tmp.mta exists", func() {
		buildProjectCmdSrc = getTestPath("mta")
		buildProjectCmdTrg = getTestPath("result")
		buildProjectCmdPlatform = "cf"
		file, err := os.Create(filepath.Join(buildProjectCmdSrc, filepath.FromSlash("Makefile_tmp.mta")))
		file.Close()
		Ω(err).Should(BeNil())
		err = buildCmd.RunE(nil, []string{})
		Ω(err).ShouldNot(BeNil())
		err = os.Remove(filepath.Join(buildProjectCmdSrc, filepath.FromSlash("Makefile_tmp.mta")))
		Ω(err).Should(BeNil())
	})
})
