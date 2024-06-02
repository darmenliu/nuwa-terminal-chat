package cmdexe

import (
	"os/exec"
)

func ExecCommandWithOutput(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func ExecCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	return cmd.Run()
}
