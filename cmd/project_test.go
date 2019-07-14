package commands

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

var _ = Describe("Project", func() {
	It("Invalid command call", func() {
		err := projectBuildCmd.RunE(nil, []string{})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(Equal(fmt.Sprintf(artifacts.UnsupportedPhaseMsg, "")))
	})
})
