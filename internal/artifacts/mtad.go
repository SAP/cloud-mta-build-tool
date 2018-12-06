package artifacts

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

// GenMtad generates an mtad.yaml file from a mta.yaml file and a platform configuration file.
func GenMtad(mtaStr *mta.MTA, ep *dir.Loc, platform string) error {
	// Create META-INF folder under the mtar folder
	metaPath, err := ep.GetMetaPath()
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed on Get Meta Path")
	}
	err = dir.CreateDirIfNotExist(metaPath)
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed, not able to create dir")
	}
	if !ep.IsDeploymentDescriptor() {
		err = ConvertTypes(*mtaStr, platform)
		if err != nil {
			return errors.Wrap(err, "mtad.yaml generation failed on type conversion")
		}
	}
	// Create readable Yaml before writing to file
	mtad, err := marshal(mtaStr)
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed on MTAD marshaling")
	}
	mtadPath, err := ep.GetMtadPath()
	if err == nil {
		// Write back the MTAD to the META-INF folder
		err = ioutil.WriteFile(mtadPath, mtad, os.ModePerm)
	}
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed")
	}
	return nil
}

// marshal - serializes the MTA into an encoded YAML document.
func marshal(in *mta.MTA) (mtads []byte, err error) {
	mtads, err = yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	return mtads, nil
}

// AdaptMtadForDeployment - remove elements from MTA that are not relevant for MTAD
// Function is used in process of deployment artifacts preparation
// SupportedPlatforms of module's build parameters indicate if module has to be deployed
// if SupportedPlatforms property defined with empty list of properties
// module will not be packed, not listed in MTAD yaml and in manifest
// if module has to be deployed we clean build parameters from module,
// as this section is not used in MTAD yaml
func AdaptMtadForDeployment(mtaStr *mta.MTA, platform string) {

	// remove build parameters from modules with defined platforms
	for _, m := range mtaStr.Modules {
		if buildops.PlatformsDefined(m) {
			m.BuildParams = map[string]interface{}{}
		}
	}

	// remove modules with no platforms defined
	for doCleaning := true; doCleaning; {
		doCleaning = false
		for i, m := range mtaStr.Modules {
			if !buildops.PlatformsDefined(m) {
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
