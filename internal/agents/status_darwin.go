// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build darwin

package agents

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	defaultStatus = "unknown"
)

func getStatus(ctx context.Context, cmd string) (string, string, string, int, error) {
	currStatus := defaultStatus
	subStatus := ""

	switch {
	case strings.HasPrefix(cmd, "brew"):
		return brewStatus(ctx, cmd)
	default:
	}

	return currStatus, subStatus, "", -1, fmt.Errorf("unable to obtain status")
}

func brewStatus(ctx context.Context, cmd string) (string, string, string, int, error) {
	currStatus := defaultStatus
	subStatus := ""

	if !strings.HasSuffix(cmd, "--json") {
		cmd += " --json"
	}

	output, exitCode, err := execute(ctx, cmd)
	if err != nil {
		return currStatus, subStatus, base64.StdEncoding.EncodeToString(output), exitCode, fmt.Errorf("%s: %w", cmd, err)
	}

	if exitCode != 0 {
		return currStatus, subStatus, base64.StdEncoding.EncodeToString(output), exitCode, fmt.Errorf("%s: %w", cmd, err)
	}

	if bytes.Contains(output, []byte(`"running": true`)) {
		currStatus = "running"
	} else {
		currStatus = "stopped"
	}

	return currStatus, subStatus, base64.StdEncoding.EncodeToString(output), exitCode, nil
}
