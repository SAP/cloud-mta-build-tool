package dir

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFsys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fsys Suite")
}
