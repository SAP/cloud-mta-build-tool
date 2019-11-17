package commands

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

var _ = Describe("Commands", func() {

	BeforeEach(func() {
		cleanupCmdTrg = getTestPath("result")
		logs.Logger = logs.NewLogger()
		Ω(dir.CreateDirIfNotExist(cleanupCmdTrg)).Should(Succeed())
	})

	AfterEach(func() {
		Ω(os.RemoveAll(cleanupCmdTrg)).Should(Succeed())
	})

	var _ = Describe("cleanup command", func() {
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
