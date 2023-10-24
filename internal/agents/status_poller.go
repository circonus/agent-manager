// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type StatusPoller struct {
	interval time.Duration
}

func NewStatusPoller() (*StatusPoller, error) {
	pi := viper.GetString(keys.StatusPollingInterval)

	i, err := time.ParseDuration(pi)
	if err != nil {
		return nil, fmt.Errorf("parsing status polling interval: %w", err)
	}

	return &StatusPoller{interval: i}, nil
}

func (p *StatusPoller) Start(ctx context.Context) {
	log.Info().Str("interval", p.interval.String()).Msg("starting poller")

	for {
		t := time.NewTimer(p.interval)
		select {
		case <-ctx.Done():
			if !t.Stop() {
				<-t.C
			}

			return
		case <-t.C:
			log.Debug().Msg("collecting agent status")

			agents, err := registration.LoadInstalledAgents()
			if err != nil {
				log.Error().Err(err).Msg("loading installed agents, run --inventory again to generate")

				continue
			}

			for _, a := range agents {
				agent, err := inventory.GetAgent(a.AgentTypeID)
				if err != nil {
					log.Error().Err(err).Str("agent_type", a.AgentTypeID).Msg("getting agent")

					continue
				}

				if agent.Status != "" {
					if err := p.submitAgentStatus(ctx, a.AgentID, agent.Status); err != nil {
						log.Warn().Err(err).Msg("submitting agent status")
					}
				}
			}
		}
	}
}

type StatusResult struct {
	Status     string     `json:"status"`
	StatusData StatusData `json:"status_data"`
}

type StatusData struct {
	SubStatus string `json:"substatus"`
	Error     string `json:"error"`
	RawResult string `json:"raw_result"`
	ExitCode  int    `json:"exit_code"`
}

func (p *StatusPoller) submitAgentStatus(ctx context.Context, agentID, cmd string) error {
	status, subStatus, statusData, exitCode, err := getStatus(ctx, cmd)
	if err != nil {
		log.Warn().Err(err).
			Str("agent_id", agentID).
			Str("status_data", statusData).
			Str("status", status).
			Str("sub_status", subStatus).
			Int("exit_code", exitCode).
			Str("cmd", cmd).
			Msg("status command failed")
	}

	result := StatusResult{
		Status: status,
		StatusData: StatusData{
			SubStatus: subStatus,
			ExitCode:  exitCode,
		},
	}

	if err != nil {
		result.StatusData.Error = err.Error()
	}

	if len(statusData) > 0 {
		result.StatusData.RawResult = statusData
	}

	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent", agentID)
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("calling actions endpoint: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response -- status: %s, body: %s", resp.Status, string(body)) //nolint:goerr113
	}

	if len(body) > 0 {
		log.Debug().RawJSON("body", body).Msg("action result response")
	}

	return nil
}
