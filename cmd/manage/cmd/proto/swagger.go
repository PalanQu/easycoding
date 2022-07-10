package proto

import (
	"easycoding/common/exec"
	"easycoding/common/workspace"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	excludeDir = "google"
	outputName = "api.swagger.json"
)

var generateCmd = &cobra.Command{
	Use:   "gen-swagger",
	Short: "generate the api/api.swagger.json files",
	Run: func(cmd *cobra.Command, args []string) {
		dryrun := getDryrun(cmd)
		commands := []string{"swagger", "mixin"}

		swaggerFiles := getSwaggerJsonFiles()
		commands = append(commands, swaggerFiles...)
		commands = append(commands,
			"--output", filepath.Join(workspace.GetWorkspace(), apiDir, outputName))
		commands = append(commands, "--quiet")
		cwd := filepath.Join(workspace.GetWorkspace(), apiDir)
		exec.ExecCommand(commands, cwd, dryrun)
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean-swagger",
	Short: "delete the api/api.swagger.json files",
	Run: func(cmd *cobra.Command, args []string) {
		dryrun := getDryrun(cmd)
		swaggerJson := filepath.Join(workspace.GetWorkspace(), apiDir, outputName)
		commands := []string{"rm", swaggerJson}
		cwd := filepath.Join(workspace.GetWorkspace(), apiDir)
		err := exec.ExecCommand(commands, cwd, dryrun)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func initSwagger() {
	protoCmd.AddCommand(generateCmd)
	protoCmd.AddCommand(cleanCmd)
}

func getSwaggerJsonFiles() []string {
	root := filepath.Join(workspace.GetWorkspace(), apiDir)
	files := []string{}
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			if errors.Is(err, filepath.SkipDir) {
				return nil
			}
			return err
		}
		if info.IsDir() && info.Name() == excludeDir {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() == outputName {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".swagger.json") {
			files = append(files, path)
		}
		return nil
	})
	return files
}
