package cmd

import (
	"easycoding/common/exec"
	"easycoding/common/workspace"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const apiPath = "api"

var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Lint proto source code",
	Run: func(cmd *cobra.Command, args []string) {
		files := args
		commands := []string{"buf", "lint"}
		for _, file := range files {
			newPath := strings.TrimPrefix(file, fmt.Sprintf("%s/", apiPath))
			commands = append(commands, "--path", newPath)
		}
		cwd := filepath.Join(workspace.GetWorkspace(), "api")
		if err := exec.ExecCommand(commands, cwd, dryrun); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func initProto() {
	rootCmd.AddCommand(protoCmd)
}
