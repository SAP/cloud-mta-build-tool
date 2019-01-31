package commands

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
)

var _ = Describe("Cmd", func() {
	var _ = Describe("Version", func() {
		It("Sanity", func() {
			version.VersionConfig = []byte(`
cli_version: 0.0.1
makefile_version: 10.5.3
`)
			out := executeAndProvideOutput(func() {
				versionCmd.Run(nil, []string{})
			})
			Ω(out).Should(Equal("0.0.1\n"))
		})
	})

})

var _ = Describe("Commands", func() {

	BeforeEach(func() {
		mtadCmdTrg = getTestPath("result")
		metaCmdTrg = getTestPath("result")
		mtarCmdTrg = getTestPath("result")
		packCmdTrg = getTestPath("result")
		buildCmdTrg = getTestPath("result")
		cleanupCmdTrg = getTestPath("result")
		logs.Logger = logs.NewLogger()
		err := os.Mkdir(mtadCmdTrg, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterEach(func() {
		os.RemoveAll(mtadCmdTrg)
	})

	var _ = Describe("cleanup command", func() {

		BeforeEach(func() {
			os.MkdirAll(getTestPath("resultClean", "mtahtml5", "mtahtml5"), os.ModePerm)
		})

		AfterEach(func() {
			os.RemoveAll(getTestPath("resultClean"))
		})

		It("Sanity", func() {
			// cleanup command used for test temp file removal
			cleanupCmdSrc = getTestPath("testdata", "mtahtml5")
			cleanupCmdTrg = getTestPath("testdata", "result")
			Ω(cleanupCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(getTestPath("testdata", "result", "mtahtml5")).ShouldNot(BeADirectory())
		})

	})

	var _ = Describe("Validate", func() {
		It("Invalid yaml path", func() {
			validateCmdSrc = getTestPath("mta1")
			Ω(validateCmd.RunE(nil, []string{})).Should(HaveOccurred())
		})
		It("Invalid descriptor", func() {
			validateCmdSrc = getTestPath("mta")
			validateCmdDesc = "x"
			Ω(validateCmd.RunE(nil, []string{})).Should(HaveOccurred())
		})
	})
})

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}
