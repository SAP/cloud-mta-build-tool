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
			Ω(ExecuteValidation(getTestPath("mta"), "dev", "semantic", "true", os.Getwd)).Should(Succeed())

		})
		It("Fails on location initialization", func() {
			Ω(ExecuteValidation("", "dev", "semantic", "true", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())

		})
		It("Fails on descriptor validation", func() {
			Ω(ExecuteValidation(getTestPath("mta"), "xx", "", "true", os.Getwd)).Should(HaveOccurred())

		})
		It("Fails on validation mode", func() {
			Ω(ExecuteValidation(getTestPath("mtahtml5"), "dev", "xx", "true", os.Getwd)).Should(HaveOccurred())

		})
		It("Fails on project validation", func() {
			Ω(ExecuteValidation(getTestPath("mtahtml5"), "dev", "", "true", os.Getwd)).Should(HaveOccurred())

		})
		It("Fails - wrong strictness indicator", func() {
			Ω(ExecuteValidation(getTestPath("mtahtml5"), "dev", "", "xxx", os.Getwd)).Should(HaveOccurred())

		})
		It("Not strict - succeeds", func() {
			Ω(ExecuteValidation(getTestPath("mta_wrong"), "dev", "schema", "false", os.Getwd)).Should(Succeed())

		})
		It("Strict - fails", func() {
			Ω(ExecuteValidation(getTestPath("mta_wrong"), "dev", "schema", "true", os.Getwd)).Should(HaveOccurred())

		})
	})
})
