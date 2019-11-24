package version

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVersion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Version Suite")
}

var _ = Describe("Version", func() {
	var config []byte
	BeforeEach(func() {
		config = VersionConfig
	})
	AfterEach(func() {
		VersionConfig = config
	})
	It("Sanity", func() {
		VersionConfig = []byte(`
cli_version: 5.2
makefile_version: 10.5.3
`)
		version, e := GetVersion()
		Ω(e).Should(Succeed())
		Ω(version, e).Should(Equal(Version{CliVersion: "5.2", MakeFile: "10.5.3"}))
	})
	It("parses the real version config successfully", func() {
		_, e := GetVersion()
		Ω(e).Should(Succeed())
	})
})
