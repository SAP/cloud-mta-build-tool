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

	It("assembleCommand - fails on missing mtad in the current location", func() {
		Ω(assembleCommand.RunE(nil, []string{})).Should(HaveOccurred())
	})

})
