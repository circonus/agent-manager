// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/rs/zerolog/log"
)

// NOTE: cmdReload may be called in two different contexts
//       1. as part of installing a new configuration (most common)
//       2. as a direct action command (least common, restart would probably be used instead)

func cmdReload(ctx context.Context, a inventory.Agent, command Command) {
	switch {
	case a.Reload == "":
		return
	case strings.ToLower(a.Reload) == RESTART:
		runCommand(ctx, a.Restart, command.ID)
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
		runCommand(ctx, a.Reload, command.ID)
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
