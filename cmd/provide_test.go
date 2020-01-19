package commands

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("provideModuleCmd", func() {
	It("Invalid command call", func() {
		Ω(provideModuleCmd.RunE(nil, []string{})).Should(HaveOccurred())
	})
})
