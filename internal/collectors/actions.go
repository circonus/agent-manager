package collectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// handle requesting actions from api, performing actions, and sending results back to api.

const (
	// action types.
	CONFIG  = "config"
	COMMAND = "command"

	STATUS_ACTIVE = "active"
	STATUS_ERROR  = "error"
)

type Actions []Action

type Action struct {
	Configs  map[string]Config `json:"configs" yaml:"configs"`
	ID       string            `json:"id" yaml:"id"` // not used yet, api only supports configs
	Type     string            `json:"type" yaml:"type"`
	Commands []Command         `json:"commands" yaml:"commands"`
}

// not OS commands, commands the agent knows (e.g. restart_collector, collector_status, etc.).
// how to restart a collector is in the collector inventory.
// collector status would be the result of running `systemctl status <collector>`.
type Command struct {
	ID        string `json:"id" yaml:"id"`
	Collector string `json:"collector" yaml:"collector"`
	Command   string `json:"command" yaml:"command"`
}

// Contest should be base64 encoded.
type Config struct {
	ID       string `json:"id" yaml:"id"`
	Path     string `json:"path" yaml:"path"`
	Contents string `json:"contents" yaml:"contents"`
}

type Result struct {
	ActionID      string `json:"action_id" yaml:"action_id"` // not used at this time, api only supports configs
	ConfigResult  `json:"config_resutl" yaml:"config_result"`
	CommandResult `json:"command_result" yaml:"command_result"`
}

// write result will be "OK" or the err received when trying to write the file.
// reload result will be empty or base64 encoded as it may be multi-line output.
type ConfigResult struct {
	ID         string     `json:"id" yaml:"id"`
	Status     string     `json:"status" yaml:"status"` // STATUS_ACTIVE or STATUS_ERROR
	ConfigData ConfigData `json:"data" yaml:"data"`
}

type ConfigData struct {
	WriteResult  string `json:"write_result" yaml:"write_result"`
	ReloadResult string `json:"reload_result" yaml:"reload_result"`
}

// Output will be base64 encoded.
type CommandResult struct {
	ID          string      `json:"id" yaml:"id"`
	Status      string      `json:"status" yaml:"status"` // active or error
	CommandData CommandData `json:"data" yaml:"data"`
}

type CommandData struct {
	Output   string `json:"output" yaml:"output"`
	Error    string `json:"error" yaml:"error"`
	ExitCode int    `json:"exit_code" yaml:"exit_code"`
}

func getActions(ctx context.Context) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "collector", "update")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
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

	var actions Actions
	if err := json.Unmarshal(body, &actions); err != nil {
		return fmt.Errorf("unmarshal body: %w", err)
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

func sendActionResult(ctx context.Context, r Result) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent", "update")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(data))
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

	var actions Actions
	if err := json.Unmarshal(body, &actions); err != nil {
		return fmt.Errorf("unmarshal body: %w", err)
	}

	return nil
}
