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

	It("content types - wrong config", func() {

		cfg := ContentTypeConfig
		ContentTypeConfig = []byte(`
content-types:
- extension: .json
  content-typex: "application/json"
`)
		_, err := GetContentTypes()
		Ω(err).Should(HaveOccurred())
		ContentTypeConfig = cfg
	})
})
