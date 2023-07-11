// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build linux || darwin || freebsd

package collectors

import (
	"context"
	"os/exec"
)

func execute(ctx context.Context, command string) ([]byte, int, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, 0, err //nolint:wrapcheck
	}

	return output, cmd.ProcessState.ExitCode(), nil
}
