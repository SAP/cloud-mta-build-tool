package buildops

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("BuildParams", func() {

	var _ = Describe("getBuildResultsPath", func() {
		var _ = DescribeTable("valid cases", func(module *mta.Module, expected string) {
			Ω(getBuildResultsPath(&fs.Loc{}, module)).Should(HaveSuffix(expected))
		},
			Entry("Implicit Build Results Path", &mta.Module{Path: "mPath"}, "mPath"),
			Entry("Explicit Build Results Path", &mta.Module{Path: "mPath", BuildParams: mta.BuildParameters{Path: "bPath"}}, "bPath"))

		var _ = Describe("invalid cases", func() {
			BeforeEach(func() {
				fs.GetWorkingDirectory = func() (string, error) {
					return "", errors.New("error")
				}
			})
			AfterEach(func() {
				fs.GetWorkingDirectory = fs.OsGetWd
			})

			It("Implicit", func() {
				module := mta.Module{Path: "mPath"}
				_, err := getBuildResultsPath(&fs.Loc{}, &module)
				Ω(err).Should(HaveOccurred())
			})
		})
	})

	var _ = DescribeTable("getRequiredTargetPath", func(requires mta.BuildRequires, module mta.Module, expected string) {
		Ω(getRequiredTargetPath(&fs.Loc{}, &module, &requires)).Should(HaveSuffix(expected))
	},
		Entry("Implicit Target Path", mta.BuildRequires{}, mta.Module{Path: "mPath"}, "mPath"),
		Entry("Explicit Target Path", mta.BuildRequires{TargetPath: "artifacts"}, mta.Module{Path: "mPath"}, filepath.Join("mPath", "artifacts")))

	var _ = Describe("ProcessRequirements", func() {

		AfterEach(func() {
			wd, _ := os.Getwd()
			os.RemoveAll(filepath.Join(wd, "testdata", "testproject", "moduleB"))
		})

		var _ = DescribeTable("Valid cases", func(artifacts []string, expectedPath string) {
			wd, _ := os.Getwd()
			ep := fs.Loc{SourcePath: filepath.Join(wd, "testdata", "testproject"), TargetPath: filepath.Join(wd, "testdata", "result")}
			require := mta.BuildRequires{
				Name:       "A",
				Artifacts:  artifacts,
				TargetPath: "./b_copied_artifacts",
			}
			mtaObj := mta.MTA{
				Modules: []*mta.Module{
					{
						Name: "A",
						Path: "ui5app",
					},
					{
						Name: "B",
						Path: "moduleB",
						BuildParams: mta.BuildParameters{
							Requires: []mta.BuildRequires{
								require,
							},
						},
					},
				},
			}
			Ω(ProcessRequirements(&ep, &mtaObj, &require, "B")).Should(Succeed())
			Ω(filepath.Join(wd, expectedPath)).Should(BeADirectory())
			Ω(filepath.Join(wd, expectedPath, "webapp", "Component.js")).Should(BeAnExistingFile())
		},
			Entry("Require All - list", []string{"*"}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All - single value", []string{"*"}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All From Parent", []string{"."}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts", "ui5app")))

		var _ = DescribeTable("Invalid cases", func(lp *fs.Loc, require mta.BuildRequires, mtaObj mta.MTA, moduleName string) {
			Ω(ProcessRequirements(lp, &mtaObj, &require, moduleName)).Should(HaveOccurred())
		},
			Entry("Module not defined",
				&fs.Loc{},
				mta.BuildRequires{Name: "A", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"},
				mta.MTA{Modules: []*mta.Module{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"C"),
			Entry("Required Module not defined",
				&fs.Loc{},
				mta.BuildRequires{Name: "C", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"},
				mta.MTA{Modules: []*mta.Module{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"B"),
			Entry("Target path - file",
				&fs.Loc{SourcePath: getTestPath("testbuildparams")},
				mta.BuildRequires{Name: "ui1", Artifacts: []string{"*"}, TargetPath: "file.txt"},
				mta.MTA{Modules: []*mta.Module{{Name: "ui1", Path: "ui1"}, {Name: "node", Path: "node"}}},
				"node"))

		var _ = Describe("More invalid cases", func() {
			var failsOnCall int
			var callsCounter int

			BeforeEach(func() {
				fs.GetWorkingDirectory = func() (string, error) {
					callsCounter++
					if callsCounter == failsOnCall {
						return "", errors.New("error")
					}
					return os.Getwd()
				}
				callsCounter = 0
			})
			AfterEach(func() {
				fs.GetWorkingDirectory = fs.OsGetWd
			})
			var _ = DescribeTable("Get source/target path fails", func(failsOn int) {
				failsOnCall = failsOn
				req := mta.BuildRequires{Name: "A", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"}
				mtaObj := mta.MTA{Modules: []*mta.Module{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}}
				Ω(ProcessRequirements(&fs.Loc{}, &mtaObj, &req, "B")).Should(HaveOccurred())
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
			lp := fs.Loc{
				SourcePath: getTestPath("testbuildparams"),
				TargetPath: getTestPath("result"),
			}
			mtaObj, _ := mta.ParseFile(&lp)
			for _, m := range mtaObj.Modules {
				if m.Name == "node" {
					for _, r := range m.BuildParams.Requires {
						ProcessRequirements(&lp, mtaObj, &r, "node")
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

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}
