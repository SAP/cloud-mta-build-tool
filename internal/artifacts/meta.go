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
func ExecuteGenMeta(source, target, desc string, extensions []string, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generating the metadata...")
	loc, err := dir.Location(source, target, desc, extensions, wdGetter)
	if err != nil {
		return errors.Wrap(err, "failed to generate metadata when initializing the location")
	}
	return executeGenMetaByLocation(loc, loc, platform, true)
}

func executeGenMetaByLocation(loc *dir.Loc, targetArtifacts dir.ITargetArtifacts, platform string, createMetaInf bool) error {
	// validate platform
	platform, err := validatePlatform(platform)
	if err != nil {
		return err
	}

	err = dir.CreateDirIfNotExist(targetArtifacts.GetMetaPath())
	if err != nil {
		return err
	}

	err = generateMeta(loc, targetArtifacts, loc.IsDeploymentDescriptor(), platform, createMetaInf)
	return err
}

// generateMeta - generate metadata artifacts
func generateMeta(loc *dir.Loc, targetArtifacts dir.ITargetArtifacts, deploymentDescriptor bool, platform string, createMetaInf bool) error {

	// parse MTA file
	m, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, genMetaMsg)
	}

	removeUndeployedModules(m, platform)
	// Generate meta info dir with required content
	err = genMetaInfo(loc, targetArtifacts, loc, deploymentDescriptor, platform, m, createMetaInf)
	return err
}

// genMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func genMetaInfo(source dir.ISourceModule, ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath, deploymentDesc bool,
	platform string, mtaStr *mta.MTA, createMetaInf bool) (rerr error) {

	if createMetaInf {
		// Set the MANIFEST.MF file
		err := setManifestDesc(source, ep, targetPathGetter, deploymentDesc, mtaStr.Modules, mtaStr.Resources)
		if err != nil {
			return errors.Wrap(err, genMetaPopulatingMsg)
		}
	}

	if !deploymentDesc {
		err := removeBuildParamsFromMta(targetPathGetter, mtaStr)
		if err != nil {
			return err
		}
	}

	err := genMtad(mtaStr, ep, deploymentDesc, platform, yaml.Marshal)
	if err != nil {
		return errors.Wrap(err, genMetaMTADMsg)
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
