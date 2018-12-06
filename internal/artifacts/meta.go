package artifacts

import (
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/platform"
	"cloud-mta-build-tool/mta"
)

// GenerateMeta - generate build metadata artifacts
func GenerateMeta(ep dir.ILoc, platform string) error {
	logs.Logger.Info("Starting Meta folder and related artifacts creation")

	// parse MTA file
	m, err := ep.ParseFile()
	if err != nil {
		return errors.Wrap(err, "Meta folder and related artifacts creation failed on MTA file parsing")
	}
	// read MTA extension file
	mExt, err := ep.ParseExtFile(platform)
	if err == nil {
		// merge MTA with extension
		mta.Merge(m, mExt)
	}

	AdaptMtadForDeployment(m, platform)
	// Generate meta info dir with required content
	err = GenMetaInfo(ep, ep.IsDeploymentDescriptor(), platform, m, []string{})
	if err != nil {
		return errors.Wrap(err, "Meta folder and related artifacts creation failed on META Info generation")
	}
	logs.Logger.Info("Meta folder and related artifacts creation finished successfully ")
	return nil
}

// GenMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func GenMetaInfo(ep dir.ITargetArtifacts, deploymentDesc bool, platform string, mtaStr *mta.MTA, modules []string) error {
	err := GenMtad(mtaStr, ep, deploymentDesc, platform)
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
		platform.ConvertTypes(mtaStr, platformCfg, platformName)
	}
	return err
}
