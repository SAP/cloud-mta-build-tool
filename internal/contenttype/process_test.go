package contenttype

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Process", func() {
	It("content types getting", func() {
		contentTypes, err := GetContentTypes()
		Ω(err).Should(Succeed())
		Ω(GetContentType(contentTypes, ".json")).Should(Equal("application/json"))
		_, err = GetContentType(contentTypes, ".unknown")
		Ω(err).Should(HaveOccurred())
	})
})
