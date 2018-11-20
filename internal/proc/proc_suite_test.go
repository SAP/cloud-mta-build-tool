package proc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Proc Suite")
}
