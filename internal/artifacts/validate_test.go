package artifacts

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	var _ = Describe("ExecuteValidation", func() {
		It("Sanity", func() {
			Ω(ExecuteValidation(getTestPath("mta"), "dev", "", os.Getwd)).Should(Succeed())

		})
		It("Fails on location initialization", func() {
			Ω(ExecuteValidation("", "dev", "", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())

		})
		It("Fails on descriptor validation", func() {
			Ω(ExecuteValidation(getTestPath("mta"), "xx", "", os.Getwd)).Should(HaveOccurred())

		})
		It("Fails on project validation", func() {
			Ω(ExecuteValidation(getTestPath("mtahtml5"), "dev", "", os.Getwd)).Should(HaveOccurred())

		})
	})
})
