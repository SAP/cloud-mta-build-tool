package artifacts

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta/mta"
)

// ExecuteGenMtad - generates MTAD from MTA
func ExecuteGenMtad(source, target, desc, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generation of the .mtad file started")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of the .mtad file failed when initializing the location")
	}

	mtaStr, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, "generation of the .mtad file failed when parsing the %v file", loc.GetMtaYamlFilename())
	}

	mtaExt, err := loc.ParseExtFile(platform)
	if err != nil {
		return errors.Wrapf(err, "generation of the .mtad file failed when parsing the %v file", loc.GetMtaExtYamlPath(platform))
	}

	mta.Merge(mtaStr, mtaExt)
	adaptMtadForDeployment(mtaStr, platform)

	err = genMtad(mtaStr, loc, loc.IsDeploymentDescriptor(), platform)
	if err != nil {
		return err
	}
	logs.Logger.Info("generation of the .mtad file finished successfully")
	return nil
}

// genMtad generates an mtad.yaml file from a mta.yaml file and a platform configuration file.
func genMtad(mtaStr *mta.MTA, ep dir.ITargetArtifacts, deploymentDesc bool, platform string) error {
	// Create META-INF folder under the mtar folder
	metaPath := ep.GetMetaPath()
	err := dir.CreateDirIfNotExist(metaPath)
	if err != nil {
		logs.Logger.Infof("the %v folder already exists", metaPath)
	}
	if !deploymentDesc {
		err = ConvertTypes(*mtaStr, platform)
		if err != nil {
			return errors.Wrapf(err, "generation of the .mtad file failed when converting types according to the %v platform", platform)
		}
	}
	// Create readable Yaml before writing to file
	mtad, err := yaml.Marshal(mtaStr)
	if err != nil {
		return errors.Wrap(err, "generation of the .mtad file failed when marshalling the .mtad object")
	}
	mtadPath := ep.GetMtadPath()
	// Write back the MTAD to the META-INF folder
	err = ioutil.WriteFile(mtadPath, mtad, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "generation of the .mtad file failed when writing")
	}
	return nil
}

// adaptMtadForDeployment - remove elements from MTA that are not relevant for MTAD
// Function is used in process of deployment artifacts preparation
// SupportedPlatforms of module's build parameters indicate if module has to be deployed
// if SupportedPlatforms property defined with empty list of properties
// module will not be packed, not listed in MTAD yaml and in manifest
// if module has to be deployed we clean build parameters from module,
// as this section is not used in MTAD yaml
func adaptMtadForDeployment(mtaStr *mta.MTA, platform string) {

	// remove build parameters from modules with defined platforms
	for _, m := range mtaStr.Modules {
		if buildops.PlatformDefined(m, platform) {
			m.BuildParams = map[string]interface{}{}
		}
	}

	// remove modules with no platforms defined
	for doCleaning := true; doCleaning; {
		doCleaning = false
		for i, m := range mtaStr.Modules {
			if !buildops.PlatformDefined(m, platform) {
				// join slices before and after removed module
				mtaStr.Modules = mtaStr.Modules[:i+copy(mtaStr.Modules[i:], mtaStr.Modules[i+1:])]
				doCleaning = true
				break
			}
		}
	}

	//TODO move to configuration
	if platform == "neo" {
		if mtaStr.Parameters == nil {
			mtaStr.Parameters = make(map[string]interface{})
		}
		if mtaStr.Parameters["hcp-deployer-version"] == nil {
			mtaStr.Parameters["hcp-deployer-version"] = "1.1.0"
		}
	}
}
