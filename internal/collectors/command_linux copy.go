//go:build linux

package collectors

import (
	"os/exec"
)

func exec(cmd string) (output []byte, code int, error) {
	cmd := exec.Command("bash", "-c", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, 0, err
	}

	return output, cmd.ProcessState.ExitCode(), nil
}