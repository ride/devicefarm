/*
Package util provides utility functions used by other packages.
*/
package util

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ErrDetached is returned if it looks like a Git repo is in a
// detached state.
var ErrDetached = errors.New("Your repo looks like it is in a detached state")

// GitBranch returns the current git branch for the given directory
func GitBranch(dir string) (branch string, err error) {
	cmd := Cmd(dir, "git rev-parse --abbrev-ref HEAD")
	out, err := cmd.Output()
	if err != nil {
		return
	}
	branch = strings.TrimSpace(string(out))
	if branch == "HEAD" {
		err = ErrDetached
	}
	return
}

// Cmd creates an exec.Cmd to run in the given directory.
func Cmd(dir string, command string) *exec.Cmd {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	return cmd
}

// CmdOutput represents the output of a bash command.
type CmdOutput struct {
	Cmd    string
	Output string
	Err    error
}

// RunAll runs all the given commands in the given directory, and
// returns a list of CmdOutputs.
func RunAll(dir string, commands ...string) []*CmdOutput {
	outputs := make([]*CmdOutput, len(commands))
	var i int
	var command string
	for i, command = range commands {
		command = strings.TrimSpace(command)
		out, err := Cmd(dir, command).Output()
		outputs[i] = &CmdOutput{command, strings.TrimSpace(string(out)), err}
		if err != nil {
			break
		}
	}
	return outputs[0 : i+1]
}

// CopyFile copies the contents of one file to another file. If the
// dst file already exists, its contents will be replaced.
func CopyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		closeError := dstFile.Close()
		if err == nil {
			err = closeError
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		// TODO: Not sure how to add test coverage for this line
		return
	}
	err = dstFile.Sync()
	return
}
