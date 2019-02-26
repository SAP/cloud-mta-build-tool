package commands

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

const (
	defaultPlatform string = "cf"
)

var assembleCmdSrc string
var assembleCmdTrg string
var assembleCmdMtarName string
var assembleCmdParallel string

func init() {
	assemblyCommand.Flags().StringVarP(&assembleCmdSrc,
		"source", "s", "", "the path to the MTA project; the current path is default")
	assemblyCommand.Flags().StringVarP(&assembleCmdTrg,
		"target", "t", "", "the path to the MBT results folder; the current path is default")
	assemblyCommand.Flags().StringVarP(&assembleCmdMtarName,
		"mtar", "m", "", "the archive name")
	assemblyCommand.Flags().StringVarP(&assembleCmdParallel,
		"parallel", "p", "false", "if true content copying will run in parallel")

}

// Generate mtar from build artifacts
var assemblyCommand = &cobra.Command{
	Use:       "assemble",
	Short:     "Assembles MTA Archive",
	Long:      "Assembles MTA Archive",
	ValidArgs: []string{"Deployment descriptor location"},
	Args:      cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := assembly(assembleCmdSrc, assembleCmdTrg, defaultPlatform, mtarCmdMtarName, assembleCmdParallel, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func assembly(source, target, platform, mtarName, copyInParallel string, getWd func() (string, error)) error {
	logs.Logger.Info("assembling the MTA project...")

	parallelCopy, err := strconv.ParseBool(copyInParallel)
	if err != nil {
		parallelCopy = false
	}
	// copy from source to target
	err = artifacts.CopyMtaContent(source, target, dir.Dep, parallelCopy, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when copying the MTA content")
	}
	// Generate meta artifacts
	err = artifacts.ExecuteGenMeta(source, target, dir.Dep, platform, false, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when generating the meta information")
	}
	// generate mtar
	err = artifacts.ExecuteGenMtar(source, target, strconv.FormatBool(target != ""), dir.Dep, mtarName, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when generating the MTA archive")
	}
	// cleanup
	err = artifacts.ExecuteCleanup(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when executing cleanup")
	}
	return nil
}
