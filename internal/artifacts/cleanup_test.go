package artifacts

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleanup", func() {

	BeforeEach(func() {
		os.MkdirAll(getTestPath("result", "mtahtml5"), os.ModePerm)
	})

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	It("Sanity", func() {
		Ω(ExecuteCleanup(getTestPath("mtahtml5"), getTestPath("result"), "dev", os.Getwd)).Should(Succeed())
		Ω(getTestPath("result", "mtahtml5")).ShouldNot(BeADirectory())
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
	It("Fails on cleanup open file", func() {
		f, err := os.Create(getTestPath("result1", "mtahtml5"))
		err = ExecuteCleanup(getTestPath("mtahtml5"), getTestPath("result1"), "dev", os.Getwd)
		Ω(err).Should(HaveOccurred())
		f.Close()
	})
})
