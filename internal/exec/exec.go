package exec

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/platform"
	"cloud-mta-build-tool/mta"
)

const (
	pathSep    = string(os.PathSeparator)
	dataZip    = pathSep + "data.zip"
	mtarSuffix = ".mtar"
)

func makeCommand(params []string) *exec.Cmd {
	if len(params) > 1 {
		return exec.Command(params[0], params[1:]...)
	}
	return exec.Command(params[0])
}

// Execute - Execute child process and wait to results
func Execute(cmdParams [][]string) error {

	for _, cp := range cmdParams {
		var cmd *exec.Cmd
		if cp[0] != "" {
			logs.Logger.Infof("Executing %s for module %s...", cp[1:], filepath.Base(cp[0]))
		} else {
			logs.Logger.Infof("Executing %s", cp[1:])
		}
		cmd = makeCommand(cp[1:])
		cmd.Dir = cp[0]

		// During the running process get the standard output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Wrapf(err, "%s cmd.StdoutPipe() error", cp[1:])
		}
		// During the running process get the standard output
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errors.Wrapf(err, "%s cmd.StderrPipe() error", cp[1:])
		}

		// Start indicator
		shutdownCh := make(chan struct{})
		go indicator(shutdownCh)

		// Execute the process immediately
		if err = cmd.Start(); err != nil {
			return errors.Wrapf(err, "%s command start error", cp[1:])
		}
		// Stream command output:
		// Creates a bufio.Scanner that will read from the pipe
		// that supplies the output written by the process.
		scanout, scanerr := scanner(stdout, stderr)

		if scanout.Err() != nil {
			return errors.Wrapf(err, "%s scanout error", cp[1:])
		}

		if scanerr.Err() != nil {
			return errors.Wrapf(err, "Reading %s stderr error", cp[1:])
		}
		// Get execution success or failure:
		if err = cmd.Wait(); err != nil {
			return errors.Wrapf(err, "Error running %s", cp[1:])
		}
		close(shutdownCh) // Signal indicator() to terminate
		logs.Logger.Infof("Finished %s", cp[1:])

	}
	return nil
}

func scanner(stdout io.ReadCloser, stderr io.ReadCloser) (*bufio.Scanner, *bufio.Scanner) {
	scanout := bufio.NewScanner(stdout)
	scanerr := bufio.NewScanner(stderr)
	// instructs the scanner to read the input by runes instead of the default by-lines.
	scanout.Split(bufio.ScanRunes)
	for scanout.Scan() {
		fmt.Print(scanout.Text())
	}
	// instructs the scanner to read the input by runes instead of the default by-lines.
	scanerr.Split(bufio.ScanRunes)
	for scanerr.Scan() {
		fmt.Print(scanerr.Text())
	}
	return scanout, scanerr
}

// Show progress when the command is executed
// and the terminal are not providing any process feedback
func indicator(shutdownCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-shutdownCh:
			return
		}
	}
}

// GenerateMeta - generate build metadata artifacts
func GenerateMeta(ep *dir.Loc) error {
	return processMta("Metadata creation", ep, []string{}, func(file []byte, args []string) error {
		// parse MTA file
		m, err := mta.Unmarshal(file)
		CleanMtaForDeployment(m)
		if err == nil {
			// Generate meta info dir with required content
			err = GenMetaInfo(ep, m, args, func(mtaStr *mta.MTA) {
				err = ConvertTypes(*mtaStr)
			})
		}
		return err
	})
}

// GenerateMtar - generate mtar archive from the build artifacts
func GenerateMtar(ep *dir.Loc) error {
	logs.Logger.Info("MTAR Generation started")
	err := processMta("MTAR generation", ep, []string{}, func(file []byte, args []string) error {
		// read MTA
		m, err := mta.Unmarshal(file)
		if err != nil {
			return errors.Wrap(err, "MTA Process failed on yaml parsing")
		}
		targetTmpDir, err := ep.GetTargetTmpDir()
		if err != nil {
			return errors.Wrap(err, "MTA Process failed on getting target temp directory")
		}
		targetDir, err := ep.GetTarget()
		if err != nil {
			return errors.Wrap(err, "MTA Process failed on getting target directory")
		}
		// archive building artifacts to mtar
		err = dir.Archive(targetTmpDir, filepath.Join(targetDir, m.ID+mtarSuffix))
		return err
	})
	if err != nil {
		return errors.Wrap(err, "MTAR Generation failed on MTA processing")
	}
	logs.Logger.Info("MTAR Generation successfully finished")
	return nil
}

// ConvertTypes - convert types to appropriate target platform types
func ConvertTypes(mtaStr mta.MTA) error {
	// Load platform configuration file
	platformCfg, err := platform.Parse(platform.PlatformConfig)
	if err == nil {
		// Modify MTAD object according to platform types
		// Todo platform should provided as command parameter
		platform.ConvertTypes(mtaStr, platformCfg, "cf")
	}
	return err
}

// process mta.yaml file
func processMta(processName string, ep *dir.Loc, args []string, process func(file []byte, args []string) error) error {
	logs.Logger.Info("Starting " + processName)
	mf, err := dir.Read(ep)
	if err == nil {
		err = process(mf, args)
		if err == nil {
			logs.Logger.Info(processName + " finish successfully ")
		}
	} else {
		err = errors.Wrap(err, "MTA file not found")
	}
	return err
}

// PackModule - pack build module artifacts
func PackModule(ep *dir.Loc, module *mta.Module, moduleName string) error {

	if !module.PlatformsDefined() {
		return nil
	}

	if ep.IsDeploymentDescriptor() {
		return copyModuleArchive(ep, module.Path, moduleName)
	}

	logs.Logger.Infof("Pack of module %v Started", moduleName)
	// Get module relative path
	moduleZipPath, err := ep.GetTargetModuleDir(moduleName)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed getting target module directory", moduleName)
	}
	logs.Logger.Info(fmt.Sprintf("Module %v will be packed and saved in folder %v", moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err = os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on making directory %v", moduleName, moduleZipPath)
	}
	// zipping the build artifacts
	logs.Logger.Infof("Starting execute zipping module %v ", moduleName)
	moduleZipFullPath := moduleZipPath + dataZip
	sourceModuleDir, err := ep.GetSourceModuleDir(module.Path)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on getting source module directory with relative path %v",
			moduleName, module.Path)
	}
	if err = dir.Archive(sourceModuleDir, moduleZipFullPath); err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on archiving", moduleName)
	}
	logs.Logger.Infof("Pack of module %v successfully finished", moduleName)
	return nil
}

// copyModuleArchive - copies module archive to temp directory
func copyModuleArchive(ep *dir.Loc, modulePath, moduleName string) error {
	logs.Logger.Infof("Copy of module %v archive Started", moduleName)
	srcModulePath, err := ep.GetSourceModuleDir(modulePath)
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive failed getting source module directory", moduleName)
	}
	moduleSrcZip := filepath.Join(srcModulePath, "data.zip")
	moduleTrgZipPath, err := ep.GetTargetModuleDir(moduleName)
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive failed getting target module directory", moduleName)
	}
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err = os.MkdirAll(moduleTrgZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive on making directory %v", moduleName, moduleTrgZipPath)
	}
	moduleTrgZip := filepath.Join(moduleTrgZipPath, "data.zip")
	err = dir.CopyFile(moduleSrcZip, filepath.Join(moduleTrgZipPath, "data.zip"))
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive failed copying %v to %v", moduleName, moduleSrcZip, moduleTrgZip)
	}
	return nil
}

// GetValidationMode - convert validation mode flag to validation process flags
func GetValidationMode(validationFlag string) (bool, bool, error) {
	switch validationFlag {
	case "":
		return true, true, nil
	case "schema":
		return true, false, nil
	case "project":
		return false, true, nil
	}
	return false, false, errors.New("wrong argument of validation mode. Expected one of [all, schema, project]")
}

// ValidateMtaYaml - Validate MTA yaml
func ValidateMtaYaml(ep *dir.Loc, validateSchema bool, validateProject bool) error {
	if validateProject || validateSchema {
		logs.Logger.Infof("Validation of %v started", ep.MtaFilename)

		// ParseFile MTA yaml content
		yamlContent, err := dir.Read(ep)

		if err != nil {
			return errors.Wrapf(err, "Validation of %v failed on reading MTA content", ep.MtaFilename)
		}
		projectPath, err := ep.GetSource()
		if err != nil {
			return errors.Wrapf(err, "Validation of %v failed on getting source", ep.MtaFilename)
		}
		// validate mta content
		issues, err := mta.Validate(yamlContent, projectPath, validateSchema, validateProject)
		if len(issues) == 0 {
			logs.Logger.Infof("Validation of %v successfully finished", ep.MtaFilename)
		} else {
			return errors.Errorf("Validation of %v failed. Issues: \n%v %s", ep.MtaFilename, issues.String(), err)
		}
	}

	return nil
}

// GetModuleAndCommands - Get module from mta.yaml and
// commands (with resolved paths) configured for the module type
func GetModuleAndCommands(ep *dir.Loc, module string) (*mta.Module, []string, error) {
	mtaObj, err := dir.ParseFile(ep)
	if err != nil {
		return nil, nil, err
	}
	// Get module respective command's to execute
	return moduleCmd(mtaObj, module)
}

// BuildModule - builds module
func BuildModule(ep *dir.Loc, moduleName string) error {

	logs.Logger.Infof("Module %v building started", moduleName)

	// Get module respective command's to execute
	module, mCmd, err := GetModuleAndCommands(ep, moduleName)
	if err != nil {
		return errors.Wrapf(err, "Module %v building failed on getting relative path and commands", moduleName)
	}

	if !ep.IsDeploymentDescriptor() {

		// Development descriptor - build includes:
		// 1. module dependencies processing
		e := processDependencies(ep, moduleName)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on processing dependencies", moduleName)
		}

		// 2. module type dependent commands execution
		modulePath, e := ep.GetSourceModuleDir(module.Path)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on getting source module directory", moduleName)
		}

		// Get module commands
		commands := cmdConverter(modulePath, mCmd)

		// Execute child-process with module respective commands
		e = Execute(commands)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on commands execution", moduleName)
		}

		// 3. Packing the modules build artifacts (include node modules)
		// into the artifactsPath dir as data zip
		e = PackModule(ep, module, moduleName)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on module's packing", moduleName)
		}
	} else if module.PlatformsDefined() {

		// Deployment descriptor
		// copy module archive to temp directory
		err = copyModuleArchive(ep, module.Path, moduleName)
		if err != nil {
			return errors.Wrapf(err, "Module %v building failed on module's archive copy", module)
		}
	}
	return nil
}

// Get commands for specific module type
func moduleCmd(mta *mta.MTA, moduleName string) (*mta.Module, []string, error) {
	for _, m := range mta.Modules {
		if m.Name == moduleName {
			commandProvider, err := builders.CommandProvider(*m)
			if err != nil {
				return nil, nil, err
			}
			return m, commandProvider.Command, nil
		}
	}
	return nil, nil, errors.Errorf("Module %v not defined in MTA", moduleName)
}

// path and commands to execute
func cmdConverter(mPath string, cmdList []string) [][]string {
	var cmd [][]string
	for i := 0; i < len(cmdList); i++ {
		cmd = append(cmd, append([]string{mPath}, strings.Split(cmdList[i], " ")...))
	}
	return cmd
}

func processDependencies(ep *dir.Loc, moduleName string) error {
	m, err := dir.ParseFile(ep)
	if err != nil {
		return err
	}
	module, err := m.GetModuleByName(moduleName)
	if err != nil {
		return err
	}
	if module.BuildParams.Requires != nil {
		for _, req := range module.BuildParams.Requires {
			e := buildops.ProcessRequirements(ep, m, &req, module.Name)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
