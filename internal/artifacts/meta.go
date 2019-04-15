package artifacts

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta/mta"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
)

// ExecuteGenMeta - generates metadata
func ExecuteGenMeta(source, target, desc, platform string, onlyModules bool, wdGetter func() (string, error)) error {
	logs.Logger.Info("generating the metadata...")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when initializing the location")
	}
	// validate platform
	platform, err = validatePlatform(platform)
	if err != nil {
		return err
	}

	err = dir.CreateDirIfNotExist(loc.GetMetaPath())
	if err != nil {
		return err
	}

	err = generateMeta(loc, loc, loc, loc.IsDeploymentDescriptor(), platform, onlyModules)
	if err != nil {
		return err
	}
	return nil
}

// generateMeta - generate metadata artifacts
func generateMeta(parser dir.IMtaParser, ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath,
	deploymentDescriptor bool, platform string, onlyModules bool) error {

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

	removeUndeployedModules(m, platform)

	// Generate meta info dir with required content
	err = genMetaInfo(ep, targetPathGetter, deploymentDescriptor, platform, m, []string{}, onlyModules)
	if err != nil {
		return err
	}
	return nil
}

// genMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func genMetaInfo(ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath, deploymentDesc bool,
	platform string, mtaStr *mta.MTA, modules []string, onlyModules bool) (rerr error) {

	// Set the MANIFEST.MF file
	err := setManifestDesc(ep, targetPathGetter, mtaStr.Modules, mtaStr.Resources, modules, onlyModules)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when populating the manifest file")
	}

	err = removeBuildParamsFromMta(targetPathGetter, mtaStr)
	if err != nil {
		return err
	}

	err = genMtad(mtaStr, ep, deploymentDesc, platform, yaml.Marshal)
	if err != nil {
		return errors.Wrap(err, "generation of metadata failed when generating the MTAD file")
	}

	return nil
}

// ConvertTypes - convert types to appropriate target platform types
func ConvertTypes(mtaStr mta.MTA, platformName string) error {
	// Load platform configuration file
	platformCfg, err := platform.Unmarshal(platform.PlatformConfig)
	if err == nil {
		// Modify MTAD object according to platform types
		platform.ConvertTypes(mtaStr, platformCfg, platformName)
	}
	return err
}
