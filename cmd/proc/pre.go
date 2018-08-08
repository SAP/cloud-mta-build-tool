package proc

import (
	"io/ioutil"
	"os"

	"encoding/json"
	"github.com/mitchellh/go-homedir"

	"cloud-mta-build-tool/cmd/constants"
	"cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta"
	"cloud-mta-build-tool/cmd/mta/models"
)

type CfgEnv struct {
	Home    string
	TmpPath string
	Id      string
	ProjDir string
}

var tmpdir string

func PreProcess() string {

	// Create temp dir

	tmpdir = dir.DefaultTempDirFunc(dir.GetPath())
	// Get Home DIR
	home, cfgPath := createCfg()

	cfg := &CfgEnv{}
	cfg.TmpPath = tmpdir
	cfg.Home = home
	cfg.ProjDir = dir.GetPath()
	// Save config struct to json
	savecfg(*cfg, cfgPath)
	return tmpdir

}

func createCfg() (string, string) {
	home, err := homedir.Dir()
	if err != nil {
		os.Exit(1)
	}
	// Config file path
	cfgPath := home + constants.PathSep + "mtacfg.json"
	// Create config dir
	dir.CreateFile(cfgPath)
	return home, cfgPath
}

func savecfg(c CfgEnv, home string) {
	jsonC, _ := json.Marshal(&c)
	ioutil.WriteFile(home, jsonC, os.ModeAppend)
}

func GetMta(wd string) models.MTA {
	// Load mta descriptor
	mtaYmlCnt := mta.Load(wd + constants.PathSep + constants.MtaYaml)
	// parse MTA
	mtaStruct, _ := mta.Parse(mtaYmlCnt)

	return mtaStruct
}

func cfgdir() (dir string) {
	hdir, _ := homedir.Dir()
	return hdir + constants.PathSep + "mtacfg.json"
}

func ReadConfig() CfgEnv {
	data, _ := ioutil.ReadFile(cfgdir())
	var cfg CfgEnv
	json.Unmarshal(data, &cfg)
	return cfg
}

func CopyModule(module string, tmpdir string) {

	// Since any command is decoupled we need to read config
	if err := dir.CopyDir(dir.GetPath()+constants.PathSep+module, tmpdir+constants.PathSep+module); err != nil {
		logs.Logger.Fatalf("Base builder: CopyDir:err: ", err)
	}

}
