package buildops_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuildops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Buildops Suite")
}
