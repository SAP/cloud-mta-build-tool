package artifacts

import (
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/platform"
	"cloud-mta-build-tool/mta"
)

// ExecuteGenMeta - generates metadata
func ExecuteGenMeta(source, target, desc, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Info("Gen META started")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "Gen META failed on location initialization")
	}
	err = generateMeta(loc, loc, loc.IsDeploymentDescriptor(), platform)
	if err != nil {
		return errors.Wrap(err, "Gen META failed")
	}
	logs.Logger.Info("Gen META successfully finished")
	return nil
}

// generateMeta - generate metadata artifacts
func generateMeta(parser dir.IMtaParser, ep dir.ITargetArtifacts, deploymentDescriptor bool, platform string) error {
	logs.Logger.Info("Starting META folder and related artifacts creation")

	// parse MTA file
	m, err := parser.ParseFile()
	if err != nil {
		return errors.Wrap(err, "META folder and related artifacts creation failed on MTA file parsing")
	}
	// read MTA extension file
	mExt, err := parser.ParseExtFile(platform)
	if err == nil {
		// merge MTA with extension
		mta.Merge(m, mExt)
	}

	adaptMtadForDeployment(m, platform)
	// Generate meta info dir with required content
	err = GenMetaInfo(ep, deploymentDescriptor, platform, m, []string{})
	if err != nil {
		return errors.Wrap(err, "META folder and related artifacts creation failed on META Info generation")
	}
	logs.Logger.Info("META folder and related artifacts creation finished successfully ")
	return nil
}

// GenMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func GenMetaInfo(ep dir.ITargetArtifacts, deploymentDesc bool, platform string, mtaStr *mta.MTA, modules []string) (rerr error) {
	err := genMtad(mtaStr, ep, deploymentDesc, platform)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed on MTAD generation")
	}
	// Create MANIFEST.MF file
	manifestPath := ep.GetManifestPath()
	file, err := dir.CreateFile(manifestPath)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed on manifest creation")
	}
	defer func() {
		errClose := file.Close()
		if errClose != nil {
			rerr = errClose
		}
	}()
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
