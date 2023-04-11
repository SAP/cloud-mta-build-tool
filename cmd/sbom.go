package commands

import (
	"os"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/spf13/cobra"
)

var projectSBomGenCmdSrc string
var projectSBomGenCmdSBOMPath string

var moduleSBomGenCmdSrc string
var moduleSBomGenCmdModules []string
var moduleSBomGenCmdAllDependencies bool
var moduleSBomGenCmdSBOMPath string

// Generate SBOM file for project
var projectSBomGenCommand = &cobra.Command{
	Use:   "sbom-gen",
	Short: "Generates SBOM for project according to configurations in the MTA development descriptor (mta.yaml)",
	Long:  "Generates SBOM for project according to configurations in the MTA development descriptor (mta.yaml)",
	Args:  cobra.MaximumNArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteProjectSBomGenerate(projectSBomGenCmdSrc, projectSBomGenCmdSBOMPath, os.Getwd)
		logError(err)
		return err
	},
	// Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Generate SBOM file for modules
var moduleSBomGenCommand = &cobra.Command{
	Use:   "module-sbom-gen",
	Short: "Generates SBOM for specified modules according to configurations in the MTA development descriptor (mta.yaml)",
	Long:  "Generates SBOM for specified modules according to configurations in the MTA development descriptor (mta.yaml)",
	Args:  cobra.MaximumNArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteModuleSBomGenerate(moduleSBomGenCmdSrc, moduleSBomGenCmdModules, moduleSBomGenCmdAllDependencies, moduleSBomGenCmdSBOMPath, os.Getwd)
		logError(err)
		return err
	},
	// Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {

	// set flags of sbom-gen command
	projectSBomGenCommand.Flags().StringVarP(&projectSBomGenCmdSrc, "source", "s", "",
		"The path of MTA project; project root path is set as default")
	projectSBomGenCommand.Flags().StringVarP(&projectSBomGenCmdSBOMPath, "sbom-file-path", "b", "",
		`The path of SBOM file, relative or absoluted; if relative path, it is relative to MTA project root; default value is <MTA project path>/<MTA project id>.bom.xml.`)

	// set flags of module-sbom-gen command
	moduleSBomGenCommand.Flags().StringVarP(&moduleSBomGenCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is set as default")
	moduleSBomGenCommand.Flags().StringSliceVarP(&moduleSBomGenCmdModules, "modules", "m", nil,
		"The names of the modules")
	moduleSBomGenCommand.Flags().BoolVarP(&moduleSBomGenCmdAllDependencies, "with-all-dependencies", "a", true,
		"Build modules including all dependencies")
	moduleSBomGenCommand.Flags().StringVarP(&moduleSBomGenCmdSBOMPath, "sbom-file-path", "b", "",
		`The path of SBOM file, relative or absoluted; if relative path, it is relative to MTA project root; default value is <MTA project path>/<MTA project id>.bom.xml.`)
}
