package proto

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"easycoding/common/exec"
	"easycoding/common/workspace"
)

const (
	apiDir       = "api"
	bufWorkFile  = "buf.work.yaml"
	apiSubfix    = "_apis"
	ignoreModels = ".buf.ignore.yaml"
)

// The schema to buf.work.yaml
type bufWork struct {
	Directories []string
}

// The schema to .buf.ignore.yaml
type bufIgnore struct {
	IgnoreModels []string `yaml:"ignore_modules"`
}

var genGoCmd = &cobra.Command{
	Use:   "gen-go",
	Short: "Generate golang protobuf, grpc, grpc gateway, swagger json files",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryrun := getDryrun(cmd)
		commands := []string{"buf", "generate"}
		cwd := filepath.Join(workspace.GetWorkspace(), apiDir)
		err := exec.ExecCommand(commands, cwd, dryrun)
		if err != nil {
			return err
		}
		delDirNames, err := getIgnoreDirs()
		if err != nil {
			return err
		}
		if len(delDirNames) == 0 {
			return nil
		}
		rmCommands := []string{"rm", "-r"}
		rmCommands = append(rmCommands, delDirNames...)

		return exec.ExecCommand(rmCommands, cwd, dryrun)
	},
}

var cleanGoCmd = &cobra.Command{
	Use:   "clean-go",
	Short: "Clean generate files",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryrun := getDryrun(cmd)
		c, err := os.ReadFile(filepath.Join(workspace.GetWorkspace(), apiDir, bufWorkFile))
		if err != nil {
			return err
		}
		bw := bufWork{}
		if err := yaml.Unmarshal(c, &bw); err != nil {
			return err
		}
		delDirNames := []string{}
		for _, d := range bw.Directories {
			delDirNames = append(delDirNames, strings.TrimSuffix(d, apiSubfix))
		}
		commands := []string{"rm", "-r"}
		commands = append(commands, delDirNames...)
		cwd := filepath.Join(workspace.GetWorkspace(), apiDir)
		exec.ExecCommand(commands, cwd, dryrun)
		return nil
	},
}

func initGo() {
	protoCmd.AddCommand(genGoCmd)
	protoCmd.AddCommand(cleanGoCmd)
}

func getIgnoreDirs() ([]string, error) {
	c, err := os.ReadFile(filepath.Join(workspace.GetWorkspace(), apiDir, ignoreModels))
	if err != nil {
		return nil, err
	}
	bw := bufIgnore{}
	if err := yaml.Unmarshal(c, &bw); err != nil {
		return nil, err
	}
	return bw.IgnoreModels, nil
}
