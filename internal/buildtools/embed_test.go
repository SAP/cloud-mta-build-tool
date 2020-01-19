package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Embed", func() {
	BeforeEach(func() {
		templatePath = ""
		Ω(os.MkdirAll("./testdata/result", os.ModePerm)).Should(Succeed())
	})

	AfterEach(func() {
		wd, _ := os.Getwd()
		Ω(os.RemoveAll(filepath.Join(wd, "testdata", "result"))).Should(Succeed())
	})

	It("sanity", func() {
		os.Args = []string{"app", "-source=./testdata/cfg.yaml", "-target=./testdata/result/cfg.txt", "-package=testpackage", "-name=Config"}
		main()
		actualContent, _ := ioutil.ReadFile("./testdata/result/cfg.txt")
		expectedContent, _ := ioutil.ReadFile("./testdata/goldenCfg.txt")
		Ω(removeSpecialSymbols(actualContent)).Should(Equal(removeSpecialSymbols(expectedContent)))
	})

	It("negative - fails on source", func() {
		os.Args = []string{"app", "-source=./testdata/cfgNotExisting.yaml", "-target=./testdata/result/cfg1.txt", "-package=testpackage", "-name=Config"}
		Ω(main).Should(Panic())
	})

	It("negative - fails on target creation", func() {
		Ω(os.Mkdir("./testdata/result/cfg2", os.ModePerm)).Should(Succeed())
		os.Args = []string{"app", "-source=./testdata/cfg.yaml", "-target=./testdata/cfg2", "-package=testpackage", "-name=Config"}
		Ω(main).Should(Panic())
	})

	It("getConf fails on empty source", func() {
		Ω(genConf("", "", "package", "var")).Should(HaveOccurred())
	})

	It("getConf fails on empty target", func() {
		Ω(genConf("./testdata/cfg.yaml", "", "package", "var")).Should(HaveOccurred())
	})

})

func removeSpecialSymbols(b []byte) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	s := string(b)
	s = strings.Replace(s, "0xd, ", "", -1)
	s = reg.ReplaceAllString(s, "")
	return s
}
