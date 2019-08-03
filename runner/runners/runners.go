package runners

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func firstLine(output string) string {
	temp := strings.Split(output, "\n")
	return temp[0]
}

func FirstLineFromCommand(path string, args []string) (string, error) {
	cmd := exec.Command(path, "--version")
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	outLine := firstLine(stdout.String())
	errLine := firstLine(stderr.String())
	if len(outLine) != 0 {
		return outLine, nil
	}
	if len(errLine) != 0 {
		return errLine, nil
	}
	return "", errors.New("No output from command")
}
