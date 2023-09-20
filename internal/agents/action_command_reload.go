// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
)

var serviceRx = regexp.MustCompile(`^[a-zA-z][a-zA-Z0-9\.\_\-]+$`)

// NOTE: cmdReload may be called in two different contexts
//       1. as part of installing a new configuration (most common)
//       2. as a direct action command (least common, restart would probably be used instead)

func cmdReload(ctx context.Context, a Agent, command Command) {
	switch {
	case a.Reload == "":
		// agent is capable of actively watching configs
		// and dynamically reloading when changes detected.
		return
	case strings.HasPrefix(strings.ToLower(a.Reload), "http"):
		// http|method|body|url -- e.g. for fluent-bit "http|post||http://localhost:2020/api/v2/reload"
		// fluent-bit -- https://docs.fluentbit.io/manual/administration/hot-reload#via-http
		parts := strings.SplitN(a.Reload, "|", 4)
		if len(parts) != 4 {
			log.Warn().Str("reload", a.Reload).Msg("invalid reload http setting")

			return
		}

		method := strings.ToUpper(parts[1])
		body := parts[2]
		rawURL := parts[3]

		respBody, err := httpReloadRequest(ctx, method, body, rawURL)
		if err != nil {
			log.Warn().Err(err).Str("reload", a.Reload).Msg("http reload failed")
		}

		if command.ID != "" {
			result := CommandResult{
				ID: command.ID,
				CommandData: CommandData{
					ExitCode: 0, // not applicable to http reloads
				},
			}
			if err != nil {
				result.CommandData.Error = err.Error()
			}

			if len(respBody) > 0 {
				result.CommandData.Output = base64.StdEncoding.EncodeToString(respBody)
			}

			if err = sendCommandResult(ctx, result); err != nil {
				log.Error().Err(err).Msg("command result")
			}
		}
	default:
		switch runtime.GOOS {
		case "darwin":
			if strings.HasPrefix(strings.ToLower(a.Reload), "brew") {
				switch runtime.GOARCH {
				case "amd64":
					brew := "/usr/local/bin/brew"
					if !doesCommandExist(brew) {
						return
					}

					log.Warn().Str("brew", brew).Msg("darwin brew reload prefix not implemented yet")
				case "arm64":
					brew := "/opt/homebrew/bin/brew"
					if !doesCommandExist(brew) {
						return
					}

					log.Warn().Str("brew", brew).Msg("darwin brew reload prefix not implemented yet")
				default:
					log.Error().Str("arch", runtime.GOARCH).Str("platform", runtime.GOOS).Msg("unsupported arch for platform")

					return
				}
			}
		case "linux":
			if strings.HasPrefix(strings.ToLower(a.Reload), "systemd") {
				_, serviceName, found := strings.Cut(a.Reload, "|")
				if !found {
					log.Error().Str("reload", a.Reload).Msg("separator not found")

					return
				}

				if serviceName == "" {
					log.Error().Str("reload", a.Reload).Msg("missing service name after separator")

					return
				}

				if !serviceRx.MatchString(serviceName) {
					log.Error().Str("reload", a.Reload).Str("rx", serviceRx.String()).Msg("invalid service name")

					return
				}

				if !isServiceEnabled(ctx, serviceName) {
					return
				}

				runCommand(ctx, "systemctl restart "+serviceName+".service", command.ID)
			}
		default:
			log.Error().Str("platform", runtime.GOOS).Msg("no reload support")
		}
	}
}

func httpReloadRequest(ctx context.Context, method, body, rawURL string) ([]byte, error) {
	_, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parsing url %s: %w", rawURL, err)
	}

	client := http.DefaultClient

	var req *http.Request

	var rerr error

	switch method {
	case http.MethodGet:
		req, rerr = http.NewRequestWithContext(ctx, method, rawURL, nil)
	case http.MethodPost, http.MethodPut:
		req, rerr = http.NewRequestWithContext(ctx, method, rawURL, strings.NewReader(body))
	default:
		return nil, fmt.Errorf("http reload, unsupported method '%s'", method) //nolint:goerr113
	}

	if rerr != nil {
		return nil, rerr //nolint:wrapcheck
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return b, fmt.Errorf("non-200 response %d - %s: %s", resp.StatusCode, resp.Status, string(b)) //nolint:goerr113
	}

	return b, nil
}

func isServiceEnabled(ctx context.Context, serviceName string) bool {
	cmd := fmt.Sprintf("systemctl is-enabled %s.service", serviceName)

	output, code, err := execute(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Str("output", string(output)).Int("exit_code", code).Str("cmd", cmd).Msg("command failed")

		return false
	}

	return true
}

func doesCommandExist(cmd string) bool {
	if _, err := os.Stat(cmd); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		log.Error().Err(err).Str("cmd", cmd).Msg("command does not exist, cannot reload agent")

		return false
	}

	return false
}
