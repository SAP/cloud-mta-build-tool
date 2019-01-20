package commands

import (
	"os"

	"cloud-mta-build-tool/internal/artifacts"
	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	defaultPlatform            string = "cf"
	defaultMtaLocation         string = ""
	defaultMtaAssemblyLocation string = ""
)

func init() {}

// Generate mtar from build artifacts
var assemblyCommand = &cobra.Command{
	Use:       "assemble",
	Short:     "Assemble MTA Archive",
	Long:      "Assemble MTA Archive",
	ValidArgs: []string{"Deployment descriptor location"},
	Args:      cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := assembly(defaultMtaLocation, defaultMtaAssemblyLocation, defaultPlatform, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}

func assembly(source, target, platform string, getWd func() (string, error)) error {
	logs.Logger.Info("assemble started")
	// copy from source to target
	err := artifacts.CopyMtaContent(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, "assemble failed when copying the mta content")
	}
	// Generate meta artifacts
	err = artifacts.ExecuteGenMeta(source, target, dir.Dep, platform, false, getWd)
	if err != nil {
		return errors.Wrap(err, "assemble failed when generating the meta info")
	}
	// generate mtar
	err = artifacts.ExecuteGenMtar(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, "assemble failed when generating the MTA archive")
	}
	err = artifacts.ExecuteCleanup(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, "assemble failed when executing cleanup")
	}
	logs.Logger.Info("assemble finished successfully")
	return nil
}
