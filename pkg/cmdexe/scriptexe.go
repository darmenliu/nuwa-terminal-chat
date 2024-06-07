package cmdexe

import "os/exec"

// ExecScript executes a shell script
func ExecScript(script string) error {
	cmd := exec.Command("bash", "-x", script)
	return cmd.Run()
}

// ExecScriptWithOutput executes a shell script and returns the output
func ExecScriptWithOutput(script string) (string, error) {
	cmd := exec.Command("bash", "-x", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
