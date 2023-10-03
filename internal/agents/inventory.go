// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	noVersion = "v0.0.0"
)

// handle requesting list of agents from api, determining if any are installed locally, and responding with agents found

// key1 is platform (GOOS) e.g. darwin, windows, linux, freebsd, etc.
// key2 is agent e.g. fluent-bit, telegraf, etc.
type Agents map[string]map[string]Agent

type Agent struct {
	ConfigFiles map[string]string `json:"config_files" yaml:"config_files"`
	Binary      string            `json:"binary"       yaml:"binary"`
	Start       string            `json:"start"        yaml:"start"`
	Stop        string            `json:"stop"         yaml:"stop"`
	Restart     string            `json:"restart"      yaml:"restart"`
	Reload      string            `json:"reload"       yaml:"reload"`
	Status      string            `json:"status"       yaml:"status"`
	Version     string            `json:"version"      yaml:"version"`
}

type InstalledAgents []InstalledAgent
type InstalledAgent struct {
	AgentTypeID string `json:"agent_type_id"`
	Version     string `json:"version"`
}

func FetchAgents(ctx context.Context) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent_type")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
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

	log.Debug().Str("resp", string(body)).Msg("response")

	agents, err := ParseAPIAgents(body)
	if err != nil {
		return fmt.Errorf("parsing api response: %w", err)
	}

	return SaveAgents(agents)
}

func LoadAgents() (Agents, error) {
	file := viper.GetString(keys.InventoryFile)
	if file == "" {
		return nil, fmt.Errorf("invalid inventory file (empty)") //nolint:goerr113
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("loading agent inventory: %w", err)
	}

	var aa Agents
	if err := yaml.Unmarshal(data, &aa); err != nil {
		return nil, fmt.Errorf("parsing agent inventory: %w", err)
	}

	return aa, nil
}

func SaveAgents(aa Agents) error {
	file := viper.GetString(keys.InventoryFile)
	if file == "" {
		return fmt.Errorf("invalid inventory file (empty)") //nolint:goerr113
	}

	data, err := yaml.Marshal(aa)
	if err != nil {
		return fmt.Errorf("marshal agent inventory: %w", err)
	}

	if err := os.WriteFile(file, data, 0600); err != nil {
		return fmt.Errorf("saving agent inventory: %w", err)
	}

	return nil
}

func CheckForAgents(ctx context.Context) error {
	aa, err := LoadAgents()
	if err != nil {
		return err
	}

	platform := getPlatform()

	gaa, ok := aa[platform]
	if !ok {
		return fmt.Errorf("no agents found for platform %s", platform) //nolint:goerr113
	}

	found := InstalledAgents{}

	for name, a := range gaa {
		if _, err := os.Stat(a.Binary); errors.Is(err, os.ErrNotExist) {
			log.Warn().Str("file", a.Binary).Msg("agent binary not found, skipping")

			continue
		}

		if viper.GetString(keys.Register) != "" {
			// if this is a registration, backup current configs
			backupConfigs(name, a.ConfigFiles)
		}

		ver, err := getAgentVersion(a.Version)
		if err != nil {
			log.Warn().Err(err).Str("agent", name).Msg("getting agent version")
		}

		log.Info().Str("agent", name).Msg("found")
		found = append(found, InstalledAgent{AgentTypeID: name, Version: ver})
	}

	if registration.IsRunningInDocker() && len(viper.GetStringSlice(keys.Agents)) > 0 {
		for _, name := range viper.GetStringSlice(keys.Agents) {
			a, ok := gaa[name]
			if !ok {
				log.Error().Str("agent", name).Msg("agent not found in inventory, skipping")

				continue
			}

			if viper.GetString(keys.Register) != "" {
				// if this is a registration, backup current configs
				backupConfigs(name, a.ConfigFiles)
			}

			log.Info().Str("agent", name).Msg("force add agent from --agents")
			found = append(found, InstalledAgent{AgentTypeID: name, Version: noVersion})
		}
	}

	if len(found) > 0 {
		// contact api and report what agents were found
		return registerAgents(ctx, found)
	}

	return nil
}

func getAgentVersion(vercmd string) (string, error) {
	if vercmd == "" {
		return noVersion, nil
	}

	cmd := exec.Command("bash", "-c", vercmd) //nolint:gosec

	output, err := cmd.CombinedOutput()
	if err != nil {
		return noVersion, fmt.Errorf("%s: %w", string(output), err)
	}

	if len(output) > 0 {
		return strings.TrimSpace(string(output)), nil
	}

	return noVersion, nil
}

func registerAgents(ctx context.Context, c InstalledAgents) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent", "manager")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal claims: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(data))
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

	log.Debug().Str("resp", string(body)).Msg("response")

	return nil
}
