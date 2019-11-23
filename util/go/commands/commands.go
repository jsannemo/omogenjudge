// Package commands contains utilities to handle executing commands.
package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func firstLine(output string) string {
	temp := strings.Split(output, "\n")
	return temp[0]
}

// FirstLineFromCommand executes a command and returns the first line from stdout, or stderr if stdout was empty.
func FirstLineFromCommand(path string, args []string) (string, error) {
	cmd := exec.Command(path, "--version")
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("could not run command %s with args %v: %v", path, args, err)
	}
	outLine := firstLine(stdout.String())
	errLine := firstLine(stderr.String())
	if len(outLine) != 0 {
		return outLine, nil
	}
	if len(errLine) != 0 {
		return errLine, nil
	}
	return "", fmt.Errorf("no output from command %s with args %v", path, args)
}
