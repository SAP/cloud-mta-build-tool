package commands

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
)

var _ = Describe("Generate commands call", func() {

	BeforeEach(func() {
		mtadGenCmdTrg = getTestPath("result")
		err := dir.CreateDirIfNotExist(mtadGenCmdTrg)
		Ω(err).Should(Succeed())
	})

	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})

	It("Generate Mtad - Sanity", func() {
		mtadGenCmdSrc = getTestPath("mtahtml5")
		mtadGenCmdPlatform = "cf"
		Ω(mtadGenCmd.RunE(nil, []string{})).Should(Succeed())
		Ω(filepath.Join(getTestPath("result"), "mtad.yaml")).Should(BeAnExistingFile())
	})

	It("Generate Mtad - Invalid source", func() {
		mtadGenCmdSrc = getTestPath("mtahtml6")
		Ω(mtadGenCmd.RunE(nil, []string{})).Should(HaveOccurred())
	})

	It("Generate Mtad - Invalid platform configuration", func() {
		mtadGenCmdSrc = getTestPath("mtahtml5")
		config := platform.PlatformConfig
		platform.PlatformConfig = []byte("wrong config")
		mtadGenCmdPlatform = "cf"
		Ω(mtadGenCmd.RunE(nil, []string{})).Should(HaveOccurred())
		platform.PlatformConfig = config
	})
})
