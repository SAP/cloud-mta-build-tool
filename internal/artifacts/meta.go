package artifacts

import (
	"cloud-mta-build-tool/internal/logs"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/platform"
	"cloud-mta-build-tool/mta"
)

// GenerateMeta - generate build metadata artifacts
func GenerateMeta(ep *dir.Loc, platform string) error {
	logs.Logger.Info("Starting Meta folder and related artifacts creation")

	// parse MTA file
	m, err := dir.ParseFile(ep)
	if err != nil {
		return errors.Wrap(err, "Meta folder and related artifacts creation failed on MTA file parsing")
	}
	// read MTA extension file
	mExt, err := dir.ParseExtFile(ep, platform)
	if err == nil {
		// merge MTA with extension
		mta.Merge(m, mExt)
	}

	AdaptMtadForDeployment(m, platform)
	// Generate meta info dir with required content
	err = GenMetaInfo(ep, platform, m, []string{}, func(mtaStr *mta.MTA, platform string) {
		err = ConvertTypes(*mtaStr, platform)
	})
	if err != nil {
		return errors.Wrap(err, "Meta folder and related artifacts creation failed on META Info generation")
	}
	logs.Logger.Info("Meta folder and related artifacts creation finished successfully ")
	return nil
}

// GenMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func GenMetaInfo(ep *dir.Loc, platform string, mtaStr *mta.MTA, modules []string, convertTypes func(mtaStr *mta.MTA, platform string)) error {
	err := GenMtad(mtaStr, ep, platform, convertTypes)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed on MTAD generation")
	}
	// Create MANIFEST.MF file
	manifestPath, err := ep.GetManifestPath()
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed on getting manifest path")
	}
	file, err := dir.CreateFile(manifestPath)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed on manifest creation")
	}
	defer file.Close()
	// Set the MANIFEST.MF file
	err = setManifetDesc(file, mtaStr.Modules, modules)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed on manifest generation")
	}

	return nil
}

// ConvertTypes - convert types to appropriate target platform types
func ConvertTypes(mtaStr mta.MTA, platformName string) error {
	// Load platform configuration file
	platformCfg, err := platform.Parse(platform.PlatformConfig)
	if err == nil {
		// Modify MTAD object according to platform types
		// Todo platform should provided as command parameter
		platform.ConvertTypes(mtaStr, platformCfg, platformName)
	}
	return err
}
