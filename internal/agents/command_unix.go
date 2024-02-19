//go:build linux || darwin || freebsd

package agents

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func execute(ctx context.Context, command string) ([]byte, int, error) {
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(c, "bash", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, cmd.ProcessState.ExitCode(), fmt.Errorf("%s: %w", command, err)
	}

	return output, cmd.ProcessState.ExitCode(), nil
}
