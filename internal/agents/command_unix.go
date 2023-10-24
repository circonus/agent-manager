// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build linux || darwin || freebsd

package agents

import (
	"context"
	"os/exec"
	"time"
)

func execute(ctx context.Context, command string) ([]byte, int, error) {
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(c, "bash", "-c", command)

	output, err := cmd.CombinedOutput()

	return output, cmd.ProcessState.ExitCode(), err
}
