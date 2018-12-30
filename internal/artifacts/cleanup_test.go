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
	It("Sanity", func() {
		file, err := os.Create(getTestPath("result", "mtahtml5", "file.txt"))
		Ω(err).Should(Succeed())
		Ω(ExecuteCleanup(getTestPath("mtahtml5"), getTestPath("result"), "dev", os.Getwd)).Should(HaveOccurred())
		Ω(getTestPath("result", "mtahtml5")).Should(BeADirectory())
		file.Close()
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
