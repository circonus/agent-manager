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

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// handle requesting actions from api, performing actions, and sending results back to api.

const (
	// action types.
	CONFIG  = "config"
	COMMAND = "command"

	STATUS_ACTIVE = "active" //nolint:revive,nosnakecase,stylecheck
	STATUS_ERROR  = "error"  //nolint:revive,nosnakecase,stylecheck
)

type Actions []Action

type Action struct {
	Configs  map[string][]Config `json:"configs"  yaml:"configs"`
	ID       string              `json:"id"       yaml:"id"` // not used yet, api only supports configs
	Type     string              `json:"type"     yaml:"type"`
	Commands []Command           `json:"commands" yaml:"commands"`
}

// not OS commands, commands the agent knows (e.g. restart_agent, agent_status, etc.).
// how to restart an agent is in the agent inventory.
// agent status would be the result of running `systemctl status <agent>`.
type Command struct {
	ID      string `json:"id"      yaml:"id"`
	Agent   string `json:"agent"   yaml:"agent"`
	Command string `json:"command" yaml:"command"`
}

// Contest should be base64 encoded.
type Config struct {
	ID       string `json:"id"       yaml:"id"`
	Path     string `json:"path"     yaml:"path"`
	Contents string `json:"contents" yaml:"contents"`
}

type Result struct {
	ActionID      string        `json:"action_id"      yaml:"action_id"` // not used at this time, api only supports configs
	ConfigResult  ConfigResult  `json:"config_result"  yaml:"config_result"`
	CommandResult CommandResult `json:"command_result" yaml:"command_result"`
}

// write result will be "OK" or the err received when trying to write the file.
// reload result will be empty or base64 encoded as it may be multi-line output.
type ConfigResult struct {
	ID         string     `json:"config_assignment_id" yaml:"config_assignment_id"`
	Status     string     `json:"status"               yaml:"status"` // STATUS_ACTIVE or STATUS_ERROR
	Info       string     `json:"info,omitempty"       yaml:"info,omitempty"`
	ConfigData ConfigData `json:"data,omitempty"       yaml:"data,omitempty"`
}

type ConfigData struct {
	WriteResult  string `json:"write_result,omitempty"  yaml:"write_result,omitempty"`
	ReloadResult string `json:"reload_result,omitempty" yaml:"reload_result,omitempty"`
}

// Output will be base64 encoded.
type CommandResult struct {
	ID          string      `json:"id"     yaml:"id"`
	Status      string      `json:"status" yaml:"status"` // active or error
	CommandData CommandData `json:"data"   yaml:"data"`
}

type CommandData struct {
	Output   string `json:"output"    yaml:"output"`
	Error    string `json:"error"     yaml:"error"`
	ExitCode int    `json:"exit_code" yaml:"exit_code"`
}

func getActions(ctx context.Context) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent", "update")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("X-Circonus-Auth-Token", token)

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

	actions, err := ParseAPIActions(body)
	if len(actions) == 0 {
		log.Debug().Msg("no actions available")

		return nil
	}

	if err != nil {
		return fmt.Errorf("parsing api actions: %w", err)
	}

	for _, action := range actions {
		switch action.Type {
		case CONFIG:
			installConfigs(ctx, action)
		case COMMAND:
			if err := runCommands(ctx, action); err != nil {
				log.Error().Err(err).Msg("running commands")
			}
		default:
			log.Warn().Str("action_type", action.Type).Msg("unknown action type, skipping")
		}
	}

	return nil
}

func sendConfigResult(ctx context.Context, r ConfigResult) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	fmt.Printf("sending config result: %v\n", string(data))

	return sendActionResult(ctx, data)
}

func sendCommandResult(ctx context.Context, r CommandResult) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	return sendActionResult(ctx, data)
}

func sendActionResult(ctx context.Context, data []byte) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent", "update")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("X-Circonus-Auth-Token", token)

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
		log.Debug().Str("body", string(body)).Msg("action result response")
	}

	return nil
}
