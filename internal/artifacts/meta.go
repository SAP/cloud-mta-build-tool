package artifacts

import (
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/platform"
	"cloud-mta-build-tool/mta"
)

// GenerateMeta - generate build metadata artifacts
func GenerateMeta(ep *dir.Loc) error {
	return processMta("Metadata creation", ep, []string{}, func(file []byte, args []string) error {
		// parse MTA file
		m, err := mta.Unmarshal(file)
		CleanMtaForDeployment(m)
		if err == nil {
			// Generate meta info dir with required content
			err = GenMetaInfo(ep, m, args, func(mtaStr *mta.MTA) {
				err = ConvertTypes(*mtaStr)
			})
		}
		return err
	})
}

// GenMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func GenMetaInfo(ep *dir.Loc, mtaStr *mta.MTA, modules []string, convertTypes func(mtaStr *mta.MTA)) error {
	err := GenMtad(mtaStr, ep, convertTypes)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	// Create MANIFEST.MF file
	manifestPath, err := ep.GetManifestPath()
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	file, err := dir.CreateFile(manifestPath)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	defer file.Close()
	// Set the MANIFEST.MF file
	err = setManifetDesc(file, mtaStr.Modules, modules)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}

	return nil
}

// ConvertTypes - convert types to appropriate target platform types
func ConvertTypes(mtaStr mta.MTA) error {
	// Load platform configuration file
	platformCfg, err := platform.Parse(platform.PlatformConfig)
	if err == nil {
		// Modify MTAD object according to platform types
		// Todo platform should provided as command parameter
		platform.ConvertTypes(mtaStr, platformCfg, "cf")
	}
	return err
}
