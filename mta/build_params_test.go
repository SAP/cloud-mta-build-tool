package mta

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
)

var _ = Describe("BuildParams", func() {

	var _ = DescribeTable("validateArtifacts", func(ep *dir.MtaLocationParameters, requiredModule *Modules, artifacts []string, match GomegaMatcher) {

		Ω(validateArtifacts(ep, requiredModule, artifacts)).Should(match)

	},
		Entry("Valid artifacts", &dir.MtaLocationParameters{}, &Modules{}, []string{"*"}, Succeed()),
		Entry("Invalid artifacts", &dir.MtaLocationParameters{}, &Modules{}, []string{"*", "."}, HaveOccurred()),
	)

	var _ = DescribeTable("getBuildResultsPath", func(module *Modules, expected string) {
		Ω(module.getBuildResultsPath(&dir.MtaLocationParameters{})).Should(HaveSuffix(expected))
	},
		Entry("Implicit Build Results Path", &Modules{Path: "mPath"}, "mPath"),
		Entry("Explicit Build Results Path", &Modules{Path: "mPath", BuildParams: buildParameters{Path: "bPath"}}, "bPath"))

	var _ = DescribeTable("getRequiredTargetPath", func(requires BuildRequires, module Modules, expected string) {
		Ω(requires.getRequiredTargetPath(&dir.MtaLocationParameters{}, &module)).Should(HaveSuffix(expected))
	},
		Entry("Implicit Target Path", BuildRequires{}, Modules{Path: "mPath"}, "mPath"),
		Entry("Explicit Target Path", BuildRequires{TargetPath: "artifacts"}, Modules{Path: "mPath"}, filepath.Join("mPath", "artifacts")))

	var _ = Describe("ProcessRequirements", func() {

		AfterEach(func() {
			wd, _ := os.Getwd()
			os.RemoveAll(filepath.Join(wd, "testdata", "testproject", "moduleB"))
		})

		var _ = DescribeTable("Valid cases", func(artifacts string, expectedPath string) {
			wd, _ := os.Getwd()
			ep := dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata", "testproject"), TargetPath: filepath.Join(wd, "testdata", "result")}
			require := BuildRequires{
				Name:       "A",
				Artifacts:  artifacts,
				TargetPath: "b_copied_artifacts",
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
			Entry("Require All - list", "[*]", filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All - single value", "*", filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All From Parent", ".", filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts", "ui5app")))

		var _ = DescribeTable("Invalid cases", func(require BuildRequires, mtaObj MTA, moduleName string) {
			Ω(require.ProcessRequirements(&dir.MtaLocationParameters{}, &mtaObj, moduleName)).Should(HaveOccurred())
		},
			Entry("Module not defined",
				BuildRequires{Name: "A", Artifacts: "[*]", TargetPath: "b_copied_artifacts"},
				MTA{Modules: []*Modules{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"C"),
			Entry("Required Module not defined",
				BuildRequires{Name: "C", Artifacts: "[*]", TargetPath: "b_copied_artifacts"},
				MTA{Modules: []*Modules{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"B"),
			Entry("Empty list of services",
				BuildRequires{Name: "A", Artifacts: "", TargetPath: "b_copied_artifacts"},
				MTA{Modules: []*Modules{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"B"),
			Entry("Wrong set of artifacts",
				BuildRequires{Name: "A", Artifacts: "[*,.]", TargetPath: "b_copied_artifacts"},
				MTA{Modules: []*Modules{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"B"))
	})

})
