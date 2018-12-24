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
	logs.Logger.Info("generation of metadata started")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when initializing location")
	}
	err = generateMeta(loc, loc, loc.IsDeploymentDescriptor(), platform)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed")
	}
	logs.Logger.Info("generation of metadata finished successfully")
	return nil
}

// generateMeta - generate metadata artifacts
func generateMeta(parser dir.IMtaParser, ep dir.ITargetArtifacts, deploymentDescriptor bool, platform string) error {

	// parse MTA file
	m, err := parser.ParseFile()
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when parsing the MTA file")
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
		return err
	}
	return nil
}

// GenMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func GenMetaInfo(ep dir.ITargetArtifacts, deploymentDesc bool, platform string, mtaStr *mta.MTA, modules []string) (rerr error) {
	err := genMtad(mtaStr, ep, deploymentDesc, platform)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when generating .mtad")
	}
	// Create MANIFEST.MF file
	manifestPath := ep.GetManifestPath()
	file, err := dir.CreateFile(manifestPath)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when creating manifest")
	}
	defer func() {
		errClose := file.Close()
		if errClose != nil {
			rerr = errClose
		}
	}()
	// Set the MANIFEST.MF file
	err = setManifestDesc(file, mtaStr.Modules, modules)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when filling manifest")
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
