package commands

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pModuleCmd", func() {
	It("Invalid command call", func() {
		Ω(pModuleCmd.RunE(nil, []string{})).Should(HaveOccurred())
	})
})
