package commands

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assembly", func() {

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})

	It("assemblyCommand - fails on missing mtad in the current location", func() {
		Ω(assemblyCommand.RunE(nil, []string{})).Should(HaveOccurred())
	})

})
