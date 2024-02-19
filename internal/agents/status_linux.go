//go:build linux

package agents

import (
	"bufio"
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
	case strings.HasPrefix(cmd, "systemctl"):
		return systemctlStatus(ctx, cmd)
	case strings.HasPrefix(cmd, "brew"):
		return brewStatus(ctx, cmd)
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

func systemctlStatus(ctx context.Context, cmd string) (string, string, string, int, error) {
	currStatus := defaultStatus
	subStatus := ""

	cmd2 := strings.Replace(cmd, "status", "show", 1)

	output, exitCode, err := execute(ctx, cmd2)
	if err != nil {
		return currStatus, subStatus, base64.StdEncoding.EncodeToString(output), exitCode, fmt.Errorf("%s: %w", cmd, err)
	}

	if exitCode != 0 {
		return currStatus, subStatus, base64.StdEncoding.EncodeToString(output), exitCode, fmt.Errorf("%s: %w", cmd, err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))

	sep := "="

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "ActiveState"):
			_, status, found := strings.Cut(line, sep)
			if found {
				currStatus = status
			}
		case strings.HasPrefix(line, "SubState"):
			_, status, found := strings.Cut(line, sep)
			if found {
				subStatus = status
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return currStatus, subStatus, "error processing command output", -1, err
	}

	output, exitCode, err = execute(ctx, cmd)
	if err != nil {
		return currStatus, subStatus, base64.StdEncoding.EncodeToString(output), exitCode, fmt.Errorf("%s: %w", cmd, err)
	}

	return currStatus, subStatus, base64.StdEncoding.EncodeToString(output), exitCode, nil
}
