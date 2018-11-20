package mta

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadFile MTA", func() {

	wd, _ := os.Getwd()

	It("Valid filename", func() {
		mta, err := ReadFile(&Loc{SourcePath: filepath.Join(wd, "testdata")})
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
	})
	It("Invalid filename", func() {
		_, err := ReadFile(&Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtax.yaml"})
		Ω(err).ShouldNot(BeNil())
	})
})
