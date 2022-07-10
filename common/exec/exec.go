package exec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	quote "github.com/kballard/go-shellquote"
)

func ExecCommand(command []string, cwd string, dryRun bool) error {
	cmd := exec.Command(command[0], command[1:]...)

	if cwd == "" {
		curr, err := os.Getwd()
		if err != nil {
			return err
		}
		cwd = curr
	}
	cmd.Dir = cwd
	color.Cyan("pushd " + cwd)
	color.Cyan("\t" + quote.Join(cmd.Args...))
	color.Cyan("popd")

	if dryRun {
		return nil
	}

	var errb bytes.Buffer
	var outb bytes.Buffer
	cmd.Stderr = &errb
	cmd.Stdout = &outb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(err.Error(), outb.String(), errb.String())
	}
	return nil
}
