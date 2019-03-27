package commands

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
	It("Invalid command call", func() {
		err := projectBuildCmd.RunE(nil, []string{})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(Equal(`the "" phase of mta project build is invalid; supported phases: "pre", "post"`))
	})
})
