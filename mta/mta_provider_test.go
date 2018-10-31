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
		ep := dir.EndPoints{SourcePath: filepath.Join(wd, "testdata")}
		mta, err := ReadMta(ep)
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
	})
	It("Invalid filename", func() {
		ep := dir.EndPoints{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtax.yaml"}
		_, err := ReadMta(ep)
		Ω(err).ShouldNot(BeNil())
	})
})
