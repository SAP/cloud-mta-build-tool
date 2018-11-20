package tpl

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/builders"
	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/proc"
	"cloud-mta-build-tool/internal/version"
	"cloud-mta-build-tool/mta"
)

const (
	makefile = "Makefile.mta"
)

type tplCfg struct {
	tplContent  []byte
	relPath     string
	preContent  []byte
	postContent []byte
	depDesc     string
}

// Make - Generate the makefile
func Make(ep *mta.Loc, mode string) error {
	tpl, err := getTplCfg(mode, ep.IsDeploymentDescriptor())
	if err != nil {
		return err
	}
	if err == nil {
		tpl.depDesc = ep.Descriptor
		// Get project working directory
		err = makeFile(ep, makefile, &tpl)
	}
	return err
}

func makeFile(ep *mta.Loc, makeFilename string, tpl *tplCfg) error {

	type api map[string]string
	// template data
	var data struct {
		File mta.MTA
		API  api
		Dep  string
	}
	// ReadFile file
	m, err := mta.ReadFile(ep)
	if err != nil {
		return errors.Wrap(err, "makeFile failed reading MTA yaml")
	}

	// Template data
	data.File = *m
	data.Dep = ep.Descriptor

	// Create maps of the template method's
	t, err := mapTpl(tpl.tplContent, tpl.preContent, tpl.postContent)
	if err != nil {
		return errors.Wrap(err, "makeFile failed mapping template")
	}
	// path for creating the file
	target, err := ep.GetTarget()
	if err != nil {
		return errors.Wrap(err, "makeFile failed getting target")
	}
	path := filepath.Join(target, tpl.relPath)
	// Create make file for the template
	mf, err := createMakeFile(path, makeFilename)
	if err != nil {
		return errors.Wrap(err, "makeFile failed on file creation")
	}
	if mf != nil {
		// Execute the template
		err = t.Execute(mf, data)

		errClose := mf.Close()
		if err != nil && errClose != nil {
			err = errors.Wrapf(err, "Makefile creation failed. Closing failed with %s", errClose)
		} else if errClose != nil {
			err = errors.Wrap(errClose, "Makefile close process failed")
		}
	}
	return err
}

//noinspection GoUnusedParameter
func mapTpl(templateContent []byte, BasePreContent []byte, BasePostContent []byte) (*template.Template, error) {
	funcMap := template.FuncMap{
		"CommandProvider": builders.CommandProvider,
		"OsCore":          proc.OsCore,
		"Version":         version.GetVersion,
	}
	fullTemplate := append(BasePreContent[:], templateContent...)
	fullTemplate = append(fullTemplate, BasePostContent...)
	fullTemplateStr := string(fullTemplate)
	// parse the template txt file
	return template.New("makeTemplate").Funcs(funcMap).Parse(fullTemplateStr)
}

// Get template (default/verbose) according to the CLI flags
func getTplCfg(mode string, isDep bool) (tplCfg, error) {
	tpl := tplCfg{}
	if (mode == "verbose") || (mode == "v") {
		if isDep {
			tpl.tplContent = makeVerboseDep
		} else {
			tpl.tplContent = makeVerbose
		}
		tpl.preContent = basePreVerbose
		tpl.postContent = basePostVerbose
	} else if mode == "" {
		tpl.tplContent = makeDefault
		tpl.preContent = basePreDefault
		tpl.postContent = basePostDefault
	} else {
		return tplCfg{}, errors.New("command is not supported")
	}
	return tpl, nil
}

// Find string in arg slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func createMakeFile(path, filename string) (file *os.File, err error) {

	fullFilename := filepath.Join(path, filename)
	var mf *os.File
	if _, err = os.Stat(fullFilename); err == nil {
		logs.Logger.Warn(fmt.Sprintf("Make file %s exists", fullFilename))
		return nil, nil
	}
	mf, err = fs.CreateFile(fullFilename)
	return mf, err
}
