package dir

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MtaContent", func() {
	var _ = Describe("ParseFile MTA", func() {

		wd, _ := os.Getwd()

		It("Valid filename", func() {
			mta, err := ParseFile(&Loc{SourcePath: filepath.Join(wd, "testdata")})
			Ω(mta).ShouldNot(BeNil())
			Ω(err).Should(BeNil())
		})
		It("Invalid filename", func() {
			_, err := ParseFile(&Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtax.yaml"})
			Ω(err).ShouldNot(BeNil())
		})
	})

})
