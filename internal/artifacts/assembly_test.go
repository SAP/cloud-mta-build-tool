package artifacts

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("Assembly", func() {

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})
	It("Sanity", func() {
		err := Assembly(getTestPath("assembly-sample"),
			getTestPath("result"), "cf", "", "?", os.Getwd)
		Ω(err).Should(Succeed())
		Ω(getTestPath("result", "com.sap.xs2.samples.javahelloworld_0.1.0.mtar")).Should(BeAnExistingFile())
	})
	It("one-level-folder", func() {
		err := Assembly(getTestPath("assembly", "one-level-folder"),
			getTestPath("result"), "cf", "", "?", os.Getwd)
		Ω(err).Should(Succeed())
		mtarFile := getTestPath("result", "proj_0.1.0.mtar")
		Ω(mtarFile).Should(BeAnExistingFile())
		compareActualAndGolden(mtarFile, "MANIFEST.MF", getTestPath("assembly", "one-level-folder", "golden.mf"))
	})
	It("non-archive-path", func() {
		err := Assembly(getTestPath("assembly", "non-archive-path"),
			getTestPath("result"), "cf", "", "?", os.Getwd)
		Ω(err).Should(Succeed())
		mtarFile := getTestPath("result", "proj_0.1.0.mtar")
		Ω(mtarFile).Should(BeAnExistingFile())
		compareActualAndGolden(mtarFile, "MANIFEST.MF", getTestPath("assembly", "non-archive-path", "golden.mf"))
	})
	var _ = DescribeTable("Fails on location initialization", func(maxCalls int) {
		calls := 0
		err := Assembly("",
			getTestPath("result"), "cf", "", "true", func() (string, error) {
				calls++
				if calls >= maxCalls {
					return "", errors.New("error")
				}
				return getTestPath("assembly-sample"), nil
			})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("failed to initialize the location when getting working directory"))
	},
		Entry("fails on CopyMtaContent", 1),
		Entry("fails on ExecuteGenMeta", 2),
		Entry("fails on ExecuteGenMtar", 3),
		Entry("fails on ExecuteCleanup", 4),
	)

})

func getFileContentFromZip(path string, filename string) ([]byte, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	for _, file := range zipFile.File {
		if strings.Contains(file.Name, filename) {
			fc, err := file.Open()
			defer fc.Close()
			if err != nil {
				return nil, err
			}
			c, err := ioutil.ReadAll(fc)
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	}
	return nil, fmt.Errorf(`file "%s" not found`, filename)
}

func compareActualAndGolden(pathToZip, filenameInZip, goldenFilePath string) {
	actualContent, err := getFileContentFromZip(pathToZip, filenameInZip)
	Ω(err).Should(Succeed())
	expectedContent, err := ioutil.ReadFile(goldenFilePath)
	Ω(err).Should(Succeed())
	Ω(removeSpecialSymbols(actualContent)).Should(Equal(removeSpecialSymbols(expectedContent)))
}

