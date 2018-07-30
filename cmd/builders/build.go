package builders

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"mbtv2/cmd/mta/metainfo"
	"mbtv2/cmd/constants"
	fs "mbtv2/cmd/fsys"
	"mbtv2/cmd/logs"
	"mbtv2/cmd/proc"
)

var (
	logger *logrus.Logger
)

// MTA archives are created in a way compatible with the JAR (Java Archive) specification.
// e.g. Archive structure - mta.mtar
//  02-20-2018   META-INF/MANIFEST.MF
//  02-20-2018   META-INF/mtad.yaml
//  02-20-2018   ui5/data.zip
// This allows reusing common tools (for creating, manipulating, signing and handling such archives).
// The deployment descriptor always contains the description of the entire application (all modules and resource declarations),
// but there may be cases in which an archive doesn’t contains all MTA modules that are defined in the descriptor.
// An example may be a module with an unsupported module type.
// Another means will be used to deploy such a module,
// while the “rest” of the MTA can be bundled into the MTA archive, and can be handled in the regular way.

func init() {

	// TODO Get env
	logger = logs.NewLogger()

}

// TODO's
// 1.  Support all types of builders
// 2.  Support MTA extension
// 3.  Support Schema versions
// 4.  logger framework - Done
// 5.  unit testing - partially for the object module
// 6.  CICD
// 7.  Build opts
// 8.  zip / war artifacts
// 9.  output mtad
// 10. provide json artifacts
// 11. log according to config
// 12. build for target platform
// 13. support json desc format
// 14. create shell script for build sequence

// BuildCfg - Functional options
type BuildCfg struct {
	Target string
	// rest of runner configuration
}

// BuildProcess manage the build process
func BuildProcess(options ...func(*BuildCfg)) (buildcfg *BuildCfg, err error) {

	// Todo - support functional options
	for _, option := range options {
		logs.Logger.Info(&option)
		if err != nil {
			// Todo support functional options
		}
	}

	// TODO get from config
	logger.Infof("Starting Build Process For Target: %s", "CF")

	// performance test
	start := time.Now()
	logger.Debugln(start)
	// Get project directory
	projdir := fs.GetPath()
	// Running pre process
	// Load and parse yml & create temp dir
	tmpDir := proc.PreProcess()

	mtaStruct := proc.GetMta(projdir)
	// Build Module types according to the manifest descriptor
	for _, mod := range mtaStruct.Modules {
		switch mod.Type {
		case "html5":

			Build(NewGruntBuilder(mod.Path, mod.Name, tmpDir), projdir, fs.DefaultTempDirFunc(projdir))
		case "nodejs", "sitecontent":
			Build(NewNPMBuilder(mod.Path, mod.Name, tmpDir), projdir, fs.DefaultTempDirFunc(projdir))
		default:
			// TODO- Use ZIP builder in this case
			logger.Info("Unknown Build type")
		}

	}
	// Generate meta info dir with required content
	metainfo.GenMetaInf(tmpDir, mtaStruct)
	// Create mtar from the building artifacts
	fs.Archive(tmpDir, projdir+constants.PathSep+mtaStruct.Id+constants.MtarSuffix, tmpDir)

	// Clean up temp folder
	os.RemoveAll(tmpDir)
	elapsed := time.Since(start)
	logger.Debugf("Execution time took %s", elapsed)

	return nil, nil
}
