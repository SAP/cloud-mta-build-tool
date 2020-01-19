package tpl

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTpl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tpl Suite")

}
