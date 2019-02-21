package artifacts

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"

	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/fs"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

type mtadLoc struct {
	path string
}

func (loc *mtadLoc) GetMtadPath() string {
	return filepath.Join(loc.path, dir.Mtad)
}

func (loc *mtadLoc) GetMetaPath() string {
	return loc.path
}

func (loc *mtadLoc) GetManifestPath() string {
	return ""
}

func (loc *mtadLoc) GetMtarDir() string {
	return ""
}

// ExecuteGenMtad - generates MTAD from MTA
func ExecuteGenMtad(source, target, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generating the MTAD file...")
	loc, err := dir.Location(source, target, dir.Dev, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of the MTAD file failed when initializing the location")
	}

	// get mta object
	mtaStr, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, `generation of the MTAD file failed when parsing the "%v" file`, loc.GetMtaYamlFilename())
	}

	// get extension object if defined
	mtaExt, err := loc.ParseExtFile(platform)
	if err != nil {
		return errors.Wrapf(err, `generation of the MTAD file failed when parsing the "%v" file`, loc.GetMtaExtYamlPath(platform))
	}

	// merge mta and extension objects
	mta.Merge(mtaStr, mtaExt)
	// init mtad object from the extended mta
	adaptMtadForDeployment(mtaStr, platform)

	return genMtad(mtaStr, &mtadLoc{target}, false, platform, yaml.Marshal)
}

func validatePlatform(platform string) error {
	if platform != "xsa" && platform != "cf" && platform != "neo" {
		return fmt.Errorf("the %s deployment platform is not supported; supported values: cf, xsa, neo", platform)
	}
	return nil
}

// genMtad generates an mtad.yaml file from a mta.yaml file and a platform configuration file.
func genMtad(mtaStr *mta.MTA, ep dir.ITargetArtifacts, deploymentDesc bool, platform string,
	marshal func(interface{}) (out []byte, err error)) error {

	// validate platform
	err := validatePlatform(platform)
	if err != nil {
		return err
	}

	// Create META-INF folder under the mtar folder
	metaPath := ep.GetMetaPath()

	// if meta folder provided, mtad will be saved in this folder, so we create it if not exists
	if metaPath != "" {
		err := dir.CreateDirIfNotExist(metaPath)
		if err != nil {
			logs.Logger.Infof(`the "%v" folder already exists`, metaPath)
		}
	}
	if !deploymentDesc {
		// convert modules types according to platform
		err := ConvertTypes(*mtaStr, platform)
		if err != nil {
			return errors.Wrapf(err,
				`generation of the MTAD file failed when converting types according to the "%v" platform`,
				platform)
		}
	}
	// Create readable Yaml before writing to file
	mtad, err := marshal(mtaStr)
	if err != nil {
		return errors.Wrap(err, "generation of the MTAD file failed when marshalling the MTAD object")
	}
	mtadPath := ep.GetMtadPath()
	// Write back the MTAD to the META-INF folder
	err = ioutil.WriteFile(mtadPath, mtad, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "generation of the MTAD file failed when writing")
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
