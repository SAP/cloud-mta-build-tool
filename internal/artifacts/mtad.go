package artifacts

import (
	"io/ioutil"
	"os"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GenMtad generates an mtad.yaml file from a mta.yaml file and a platform configuration file.
func GenMtad(mtaStr *mta.MTA, ep *dir.Loc, convertTypes func(mtaStr *mta.MTA)) error {
	// Create META-INF folder under the mtar folder
	metaPath, err := ep.GetMetaPath()
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed")
	}
	err = dir.CreateDirIfNotExist(metaPath)
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed, not able to create dir")
	}
	if !ep.IsDeploymentDescriptor() {
		convertTypes(mtaStr)
	}
	// Create readable Yaml before writing to file
	mtad, err := marshal(mtaStr)
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed")
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

// CleanMtaForDeployment - remove elements from MTAR that are not relevant for MTAD
// Function is used in process of deployment artifacts preparation
// SupportedPlatforms of module's build parameters indicate if module has to be deployed
// if SupportedPlatforms property defined with empty list of properties
// module will not be packed, not listed in MTAD yaml and in manifest
// if module has to be deployed we clean build parameters from module,
// as this section is not used in MTAD yaml
func CleanMtaForDeployment(mtaStr *mta.MTA) {
	deleted := 0
	for i, m := range mtaStr.Modules {
		j := i - deleted
		if !m.PlatformsDefined() {
			// remove modules with no platforms defined
			mtaStr.Modules = mtaStr.Modules[:j+copy(mtaStr.Modules[j:], mtaStr.Modules[j+1:])]
			deleted++
		} else {
			// remove build parameters
			m.BuildParams = mta.BuildParameters{}
		}
	}
}
