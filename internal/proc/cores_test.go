package proc

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Proc", func() {
	It("OsCore", func() {
		Ω(OsCore().MAKEFLAGS).Should(Equal("MAKEFLAGS += -j"))
	})
})
