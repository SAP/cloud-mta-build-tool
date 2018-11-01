package mta_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMta(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mta Suite")
}
