package tpl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"cloud-mta-build-tool/internal/version"
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"

	"cloud-mta-build-tool/internal/builders"
	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/proc"
)

const (
	makefile        = "Makefile.mta"
	basePreVerbose  = "base_pre_verbose.txt"
	basePostVerbose = "base_post_verbose.txt"
	basePreDefault  = "base_pre_default.txt"
	basePostDefault = "base_post_default.txt"
	makeDefaultTpl  = "make_default.txt"
	makeVerboseTpl  = "make_verbose.txt"
	makeDeployTpl   = "make_deployment.txt"
)

type tplCfg struct {
	tplName string
	relPath string
	pre     string
	post    string
	depDesc string
}

// Make - Generate the makefile
func Make(ep *fs.MtaLocationParameters, mode string) error {
	tpl, err := makeMode(mode)
	if err != nil {
		return err
	}
	if ep.Descriptor == "dep" {
		tpl.tplName = makeDeployTpl
	}
	if err == nil {
		tpl.depDesc = ep.Descriptor
		// Get project working directory
		err = makeFile(ep, makefile, &tpl)
	}
	return err
}

func makeFile(ep *fs.MtaLocationParameters, makeFilename string, tpl *tplCfg) error {

	type API map[string]string
	// template data
	var data struct {
		File mta.MTA
		API  API
		Dep  string
	}
	// Read file
	m, err := mta.ReadMta(ep)
	if err != nil {
		return errors.Wrap(err, "makeFile failed reading MTA yaml")
	}

	// Template data
	data.File = *m
	data.Dep = ep.Descriptor

	// Create maps of the template method's
	t, err := mapTpl(tpl.tplName, tpl.pre, tpl.post, data.Dep != "dev")
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
	makeFile, err := createMakeFile(path, makeFilename)
	if err != nil {
		return errors.Wrap(err, "makeFile failed on file creation")
	}
	if makeFile != nil {
		// Execute the template
		err = t.Execute(makeFile, data)

		errClose := makeFile.Close()
		if err != nil && errClose != nil {
			err = errors.Wrapf(err, "Makefile creation failed. Closing failed with %s", errClose)
		} else if errClose != nil {
			err = errors.Wrap(errClose, "Makefile сlosing failed")
		}
	}
	return err
}

func mapTpl(templateName string, BasePre string, BasePost string, isDeployment bool) (*template.Template, error) {
	funcMap := template.FuncMap{
		"CommandProvider": builders.CommandProvider,
		"OsCore":          proc.OsCore,
		"Version":         version.GetVersion,
	}
	// Get the path of the template source code
	_, file, _, _ := runtime.Caller(0)
	makeVerbPath := filepath.Join(filepath.Dir(file), templateName)
	prePath := filepath.Join(filepath.Dir(file), BasePre)
	postPath := filepath.Join(filepath.Dir(file), BasePost)
	// parse the template txt file
	return template.New(templateName).Funcs(funcMap).ParseFiles(makeVerbPath, prePath, postPath)
}

// Get template (default/verbose) according to the CLI flags
func makeMode(mode string) (tplCfg, error) {
	tpl := tplCfg{}
	if (mode == "verbose") || (mode == "v") {
		tpl.tplName = makeVerboseTpl
		tpl.pre = basePreVerbose
		tpl.post = basePostVerbose
	} else if mode == "" {
		tpl.tplName = makeDefaultTpl
		tpl.pre = basePreDefault
		tpl.post = basePostDefault
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
	} else {
		mf, err = fs.CreateFile(fullFilename)
	}
	return mf, err
}
