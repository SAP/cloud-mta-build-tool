package artifacts

import (
	"errors"
	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleanup", func() {

	BeforeEach(func() {
		dir.CreateDirIfNotExist(getTestPath("result", ".mtahtml5_mta_build_tmp"))
	})

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	It("Sanity", func() {
		Ω(ExecuteCleanup(getTestPath("mtahtml5"), getResultPath(), "dev", os.Getwd)).Should(Succeed())
		Ω(getTestPath("result", ".mtahtml5_mta_build_tmp")).ShouldNot(BeADirectory())
	})
	It("Fails on location initialization", func() {
		err := ExecuteCleanup("", getTestPath("result"), "dev", func() (string, error) {
			return "", errors.New("err")
		})
		checkError(err, cleanupFailedOnLocMsg)
	})
})
