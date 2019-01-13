package exec_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/logs"
)

func TestExec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Exec Suite")
}

var _ = BeforeSuite(func() {
	logs.NewLogger()
})
