package artifacts

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
	"github.com/SAP/cloud-mta/mta"
)

// ExecuteGenMeta - generates metadata
func ExecuteGenMeta(source, target, desc string, extensions []string, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generating the metadata...")
	loc, err := dir.Location(source, target, desc, extensions, wdGetter)
	if err != nil {
		return errors.Wrap(err, "failed to generate metadata when initializing the location")
	}
	return executeGenMetaByLocation(loc, loc, platform, true, false)
}

func executeGenMetaByLocation(loc *dir.Loc, targetArtifacts dir.ITargetArtifacts, platform string, createMetaInf bool, skipValidation bool) error {
	// validate platform
	platform, err := validatePlatform(platform)
	if err != nil {
		return err
	}

	err = dir.CreateDirIfNotExist(targetArtifacts.GetMetaPath())
	if err != nil {
		return err
	}

	err = generateMeta(loc, targetArtifacts, loc.IsDeploymentDescriptor(), platform, createMetaInf, skipValidation)
	return err
}

// generateMeta - generate metadata artifacts
func generateMeta(loc *dir.Loc, targetArtifacts dir.ITargetArtifacts, deploymentDescriptor bool, platform string, createMetaInf bool, skipValidation bool) error {

	// parse MTA file
	m, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, genMetaMsg)
	}

	removeUndeployedModules(m, platform)

	setPlatformSpecificParameters(m, platform)

	// Generate meta info dir with required content
	err = genMetaInfo(loc, targetArtifacts, loc, deploymentDescriptor, platform, m, createMetaInf, skipValidation)
	return err
}

// genMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func genMetaInfo(source dir.ISourceModule, ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath, deploymentDesc bool,
	platform string, mtaStr *mta.MTA, createMetaInf bool, skipValidation bool) (rerr error) {

	if createMetaInf {
		// Set the MANIFEST.MF file
		err := setManifestDesc(source, ep, targetPathGetter, deploymentDesc, mtaStr.Modules, mtaStr.Resources)
		if err != nil {
			return errors.Wrap(err, genMetaPopulatingMsg)
		}
	}

	if !deploymentDesc {
		err := removeBuildParamsFromMta(targetPathGetter, mtaStr, skipValidation)
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

// ExecuteMerge merges mta.yaml and MTA extension descriptors and writes the result to a file with the given name
func ExecuteMerge(source, target string, extensions []string, name string, wdGetter func() (string, error)) error {
	logs.Logger.Info(mergeInfoMsg)

	if name == "" {
		return fmt.Errorf(mergeNameRequiredMsg)
	}
	loc, err := dir.Location(source, target, dir.Dev, extensions, wdGetter)
	if err != nil {
		return err
	}
	m, err := loc.ParseFile()
	if err != nil {
		return err
	}
	merged, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	mtaPath := filepath.Join(target, name)
	// Check the file doesn't already exist
	if _, err = os.Stat(mtaPath); err == nil {
		return fmt.Errorf(mergeFailedOnFileCreationMsg, mtaPath)
	}
	err = dir.CreateDirIfNotExist(filepath.Clean(target))
	if err != nil {
		return err
	}
	// Write the mta file to the selected folder
	err = ioutil.WriteFile(mtaPath, merged, os.ModePerm)
	return err
}
