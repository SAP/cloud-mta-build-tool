package commands

import (
	"testing"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuilders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builders Suite")
}

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})
