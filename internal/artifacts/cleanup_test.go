package artifacts

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleanup", func() {

	BeforeEach(func() {
		os.MkdirAll(getTestPath("result", ".mtahtml5_mta_build_tmp"), os.ModePerm)
	})

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	It("Sanity", func() {
		Ω(ExecuteCleanup(getTestPath("mtahtml5"), getResultPath(), "dev", os.Getwd)).Should(Succeed())
		Ω(getTestPath("result", ".mtahtml5_mta_build_tmp")).ShouldNot(BeADirectory())
	})
	It("Fails on location initialization", func() {
		Ω(ExecuteCleanup("", getTestPath("result"), "dev", func() (string, error) {
			return "", errors.New("err")
		})).Should(HaveOccurred())
	})
})

var _ = Describe("Cleanup", func() {
	BeforeEach(func() {
		os.MkdirAll(getTestPath("result1"), os.ModePerm)
	})
	AfterEach(func() {
		os.RemoveAll(getTestPath("result1"))
	})
})
