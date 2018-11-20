package mta

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("BuildParams", func() {

	var _ = Describe("getBuildResultsPath", func() {
		var _ = DescribeTable("valid cases", func(module *Modules, expected string) {
			Ω(module.getBuildResultsPath(&MtaLocationParameters{})).Should(HaveSuffix(expected))
		},
			Entry("Implicit Build Results Path", &Modules{Path: "mPath"}, "mPath"),
			Entry("Explicit Build Results Path", &Modules{Path: "mPath", BuildParams: buildParameters{Path: "bPath"}}, "bPath"))

		var _ = Describe("invalid cases", func() {
			BeforeEach(func() {
				dir.GetWorkingDirectory = func() (string, error) {
					return "", errors.New("error")
				}
			})
			AfterEach(func() {
				dir.GetWorkingDirectory = dir.OsGetWd
			})

			It("Implicit", func() {
				module := Modules{Path: "mPath"}
				_, err := module.getBuildResultsPath(&MtaLocationParameters{})
				Ω(err).Should(HaveOccurred())
			})
		})
	})

	var _ = DescribeTable("getRequiredTargetPath", func(requires BuildRequires, module Modules, expected string) {
		Ω(requires.getRequiredTargetPath(&MtaLocationParameters{}, &module)).Should(HaveSuffix(expected))
	},
		Entry("Implicit Target Path", BuildRequires{}, Modules{Path: "mPath"}, "mPath"),
		Entry("Explicit Target Path", BuildRequires{TargetPath: "artifacts"}, Modules{Path: "mPath"}, filepath.Join("mPath", "artifacts")))

	var _ = Describe("ProcessRequirements", func() {

		AfterEach(func() {
			wd, _ := os.Getwd()
			os.RemoveAll(filepath.Join(wd, "testdata", "testproject", "moduleB"))
		})

		var _ = DescribeTable("Valid cases", func(artifacts []string, expectedPath string) {
			wd, _ := os.Getwd()
			ep := MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata", "testproject"), TargetPath: filepath.Join(wd, "testdata", "result")}
			require := BuildRequires{
				Name:       "A",
				Artifacts:  artifacts,
				TargetPath: "./b_copied_artifacts",
			}
			mtaObj := MTA{
				Modules: []*Modules{
					{
						Name: "A",
						Path: "ui5app",
					},
					{
						Name: "B",
						Path: "moduleB",
						BuildParams: buildParameters{
							Requires: []BuildRequires{
								require,
							},
						},
					},
				},
			}
			Ω(require.ProcessRequirements(&ep, &mtaObj, "B")).Should(Succeed())
			Ω(filepath.Join(wd, expectedPath)).Should(BeADirectory())
			Ω(filepath.Join(wd, expectedPath, "webapp", "Component.js")).Should(BeAnExistingFile())
		},
			Entry("Require All - list", []string{"*"}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All - single value", []string{"*"}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All From Parent", []string{"."}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts", "ui5app")))

		var _ = DescribeTable("Invalid cases", func(lp *MtaLocationParameters, require BuildRequires, mtaObj MTA, moduleName string) {
			Ω(require.ProcessRequirements(lp, &mtaObj, moduleName)).Should(HaveOccurred())
		},
			Entry("Module not defined",
				&MtaLocationParameters{},
				BuildRequires{Name: "A", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"},
				MTA{Modules: []*Modules{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"C"),
			Entry("Required Module not defined",
				&MtaLocationParameters{},
				BuildRequires{Name: "C", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"},
				MTA{Modules: []*Modules{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"B"),
			Entry("Target path - file",
				&MtaLocationParameters{SourcePath: getTestPath("testbuildparams")},
				BuildRequires{Name: "ui1", Artifacts: []string{"*"}, TargetPath: "file.txt"},
				MTA{Modules: []*Modules{{Name: "ui1", Path: "ui1"}, {Name: "node", Path: "node"}}},
				"node"))

		var _ = Describe("More invalid cases", func() {
			var failsOnCall int
			var callsCounter int

			BeforeEach(func() {
				dir.GetWorkingDirectory = func() (string, error) {
					callsCounter++
					if callsCounter == failsOnCall {
						return "", errors.New("error")
					}
					return os.Getwd()
				}
				callsCounter = 0
			})
			AfterEach(func() {
				dir.GetWorkingDirectory = dir.OsGetWd
			})
			var _ = DescribeTable("Get source/target path fails", func(failsOn int) {
				failsOnCall = failsOn
				req := BuildRequires{Name: "A", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"}
				mtaObj := MTA{Modules: []*Modules{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}}
				Ω(req.ProcessRequirements(&MtaLocationParameters{}, &mtaObj, "B")).Should(HaveOccurred())
			},
				Entry("source", 1),
				Entry("target", 2))

		})
	})

	var _ = Describe("Process complex list of requirements", func() {
		AfterEach(func() {
			os.RemoveAll(getTestPath("testbuildparams", "node", "existingfolder", "deepfolder"))
			os.RemoveAll(getTestPath("testbuildparams", "node", "newfolder"))
		})

		It("", func() {
			lp := MtaLocationParameters{
				SourcePath: getTestPath("testbuildparams"),
				TargetPath: getTestPath("result"),
			}
			mtaObj, _ := ReadMta(&lp)
			for _, m := range mtaObj.Modules {
				if m.Name == "node" {
					for _, r := range m.BuildParams.Requires {
						r.ProcessRequirements(&lp, mtaObj, "node")
					}
				}
			}
			// ["*"] => "newfolder"
			Ω(getTestPath("testbuildparams", "node", "newfolder", "webapp")).Should(BeADirectory())
			// ["deep/folder/inui2/anotherfile.txt"] => "existingfolder/deepfolder"
			Ω(getTestPath("testbuildparams", "node", "existingfolder", "deepfolder", "anotherfile.txt")).Should(BeAnExistingFile())
			// ["./deep/*/inui2/another*"] => "./existingfolder/deepfolder"
			Ω(getTestPath("testbuildparams", "node", "existingfolder", "deepfolder", "anotherfile2.txt")).Should(BeAnExistingFile())
			// ["deep/folder/inui2/somefile.txt", "*/folder/"] =>  "newfolder/newdeepfolder"
			Ω(getTestPath("testbuildparams", "node", "newfolder", "newdeepfolder", "folder")).Should(BeADirectory())
		})

	})
})
