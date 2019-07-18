package tpl

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/proc"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta/mta"
)

type tplCfg struct {
	tplContent  []byte
	relPath     string
	preContent  []byte
	postContent []byte
	depDesc     string
}

// ExecuteMake - generate makefile
func ExecuteMake(source, target, name, mode string, wdGetter func() (string, error)) error {
	logs.Logger.Infof(`generating the "%s" file...`, name)
	loc, err := dir.Location(source, target, dir.Dev, wdGetter)
	if err != nil {
		return errors.Wrapf(err, genFailedOnInitLocMsg, name)
	}
	err = genMakefile(loc, loc, loc, name, mode)
	if err != nil {
		return err
	}
	logs.Logger.Info("done")
	return nil
}

// genMakefile - Generate the makefile
func genMakefile(mtaParser dir.IMtaParser, loc dir.ITargetPath, desc dir.IDescriptor, makeFilename, mode string) error {
	tpl, err := getTplCfg(mode, desc.IsDeploymentDescriptor())
	if err != nil {
		return err
	}
	if err == nil {
		tpl.depDesc = desc.GetDescriptor()
		// Get project working directory
		err = makeFile(mtaParser, loc, makeFilename, &tpl)
	}
	return err
}

type templateData struct {
	File mta.MTA
}

// ConvertToShellArgument wraps a string in quotation marks if necessary and escapes necessary characters in it,
// so it can be used as a shell argument or flag in the makefile
func (data templateData) ConvertToShellArgument(s string) string {
	return shellquote.Join(s)
}

// makeFile - generate makefile form templates
func makeFile(mtaParser dir.IMtaParser, loc dir.ITargetPath, makeFilename string, tpl *tplCfg) (e error) {

	// template data
	data := templateData{}

	err := dir.CreateDirIfNotExist(loc.GetTarget())
	if err != nil {
		return err
	}

	// ParseFile file
	m, err := mtaParser.ParseFile()
	if err != nil {
		return errors.Wrapf(err, genFailedOnReadMsg, makeFilename)
	}

	// Template data
	data.File = *m

	// Create maps of the template method's
	t, err := mapTpl(tpl.tplContent, tpl.preContent, tpl.postContent)
	if err != nil {
		return errors.Wrapf(err, genFailedOnTmplMapMsg, makeFilename)
	}
	// path for creating the file
	target := loc.GetTarget()

	path := filepath.Join(target, tpl.relPath)
	// Create genMakefile file for the template
	mf, err := createMakeFile(path, makeFilename)
	defer func() {
		e = dir.CloseFile(mf, e)
	}()
	if err != nil {
		return err
	}
	if mf != nil {
		// Execute the template
		err = t.Execute(mf, data)
	}
	return err
}

//noinspection GoUnusedParameter
func mapTpl(templateContent []byte, BasePreContent []byte, BasePostContent []byte) (*template.Template, error) {
	funcMap := template.FuncMap{
		"CommandProvider": func(modules mta.Module) (commands.CommandList, error) {
			cmds, _, err := commands.CommandProvider(modules)
			return cmds, err
		},
		"OsCore":  proc.OsCore,
		"Version": version.GetVersion,
	}
	fullTemplate := append(baseArgs, BasePreContent...)
	fullTemplate = append(fullTemplate, templateContent...)
	fullTemplate = append(fullTemplate, BasePostContent...)
	fullTemplateStr := string(fullTemplate)
	// parse the template txt file
	return template.New("makeTemplate").Funcs(funcMap).Parse(fullTemplateStr)
}

// Get template (default/verbose) according to the CLI flags
func getTplCfg(mode string, isDep bool) (tplCfg, error) {
	tpl := tplCfg{}
	if (mode == "verbose") || (mode == "v") {
		tpl.tplContent = makeVerbose
		tpl.preContent = basePreVerbose
		tpl.postContent = basePost
	} else if mode == "" {
		tpl.tplContent = makeDefault
		tpl.preContent = basePreDefault
		tpl.postContent = basePost
	} else {
		return tplCfg{}, fmt.Errorf(cmdNotSupportedMsg, mode)
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
		return nil, fmt.Errorf(genFailedOnFileCreationMsg, filename, fullFilename)
	}
	mf, err = dir.CreateFile(fullFilename)
	return mf, err
}
