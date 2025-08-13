package artifacts

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
)

var _ = Describe("Assembly", func() {

	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})
	It("Sanity", func() {
		err := Assembly(getTestPath("assembly-sample"), "",
			getTestPath("result"), nil, "cf", "", "?", os.Getwd)
		Ω(err).Should(Succeed())
		Ω(getTestPath("result", "com.sap.xs2.samples.javahelloworld_0.1.0.mtar")).Should(BeAnExistingFile())
	})
	It("path variations", func() {
		err := Assembly(getTestPath("assembly"), "",
			getTestPath("result"), nil, "cf", "", "?", os.Getwd)
		Ω(err).Should(Succeed())
		mtarFile := getTestPath("result", "proj_0.1.0.mtar")
		Ω(mtarFile).Should(BeAnExistingFile())
		actualContent, err := getFileContentFromZip(mtarFile, "MANIFEST.MF")
		Ω(err).Should(Succeed())
		expectedContent := getFileContentWithCliVersion(getTestPath("assembly", "golden.mf"))
		Ω(removeSpecialSymbols(actualContent)).Should(Equal(removeSpecialSymbols([]byte(expectedContent))))
	})

	var _ = DescribeTable("Fails on location initialization", func(maxCalls int) {
		calls := 0
		err := Assembly("", "",
			getTestPath("result"), nil, "cf", "", "true", func() (string, error) {
				calls++
				if calls >= maxCalls {
					return "", errors.New("error")
				}
				return getTestPath("assembly-sample"), nil
			})
		checkError(err, dir.InitLocFailedOnWorkDirMsg)
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
	defer zipFile.Close()
	var fc io.ReadCloser
	for _, file := range zipFile.File {
		if strings.Contains(file.Name, filename) {
			fc, err = file.Open()
			if err != nil {
				return nil, err
			}
			defer fc.Close() // If we got here there won't be another loop iteration
			return ioutil.ReadAll(fc)
		}
	}
	return nil, fmt.Errorf(`file "%s" not found`, filename)
}
