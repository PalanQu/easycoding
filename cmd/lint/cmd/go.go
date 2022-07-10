package cmd

import (
	"easycoding/common/exec"
	"easycoding/common/workspace"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Lint go source code",
	Long: `Use golangci-lint to lint golang source code, refactor
	golangci-lint because of typechecking error: named files must all be in
	one directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		files := args
		packages := []string{}
		for _, file := range files {
			parts := strings.Split(file, "/")
			packageName := strings.Join(parts[0:len(parts)-1], "/")
			if !contains(packages, packageName) {
				packages = append(packages, packageName)
			}
		}

		commands := []string{"golangci-lint", "run"}
		commands = append(commands, packages...)

		if err := exec.ExecCommand(commands, workspace.GetWorkspace(), dryrun); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func initGo() {
	rootCmd.AddCommand(goCmd)
}

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}
