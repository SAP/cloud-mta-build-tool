package mta

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read MTA", func() {

	wd, _ := os.Getwd()

	It("Valid filename", func() {
		mta, err := ReadMta(&dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata")})
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
	})
	It("Invalid filename", func() {
		_, err := ReadMta(&dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtax.yaml"})
		Ω(err).ShouldNot(BeNil())
	})
})
