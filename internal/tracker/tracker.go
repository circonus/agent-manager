// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package tracker

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Tracker struct {
	AgentID      string `json:"agent_id"      yaml:"agent_id"`
	AssignmentID string `json:"assignment_id" yaml:"assignment_id"`
	S            string `json:"s"             yaml:"s"`
	D            string `json:"d"             yaml:"d"`
	Modified     bool   `json:"modified"      yaml:"modified"`
}

func VerifyConfig(ctx context.Context, agentName, cfgFile string) error {
	trackerFile, err := getTrackerFile(agentName, cfgFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	t, err := loadTracker(trackerFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Warn().
				Err(err).
				Str("file", cfgFile).
				Str("agent", agentName).
				Msg("no config to track")

			return nil
		}

		return err //nolint:wrapcheck
	}

	if t.Modified {
		return nil // has already been detected and reported, short-circuit to not report over and over
	}

	if t.AgentID == "" || t.AssignmentID == "" || t.S == "" {
		return fmt.Errorf("no current tracking information available") //nolint:goerr113
	}

	s, err := generateChecksum(cfgFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if s != t.S {
		if err := UpdateAssignmentStatus(ctx, t); err != nil {
			return err
		}

		log.Warn().Str("curr", s).Str("orig", t.S).Str("id", t.AssignmentID).Msg("file modified")

		t.Modified = true
		if err := saveTracker(trackerFile, t); err != nil {
			return err //nolint:wrapcheck
		}
	}

	return nil
}

func UpdateAssignmentStatus(ctx context.Context, t *Tracker) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(
		viper.GetString(keys.APIURL),
		"agent",
		t.AgentID,
		"config_assignment",
		t.AssignmentID)
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	status := []byte(`{"status":"modified"}`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewReader(status))
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

	if resp.StatusCode == http.StatusUnauthorized {
		if err := registration.RefreshRegistration(ctx); err != nil { //nolint:govet
			return fmt.Errorf("new token: %w", err)
		}

		return fmt.Errorf("token expired, refreshed") //nolint:goerr113
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response -- status: %s, body: %s", resp.Status, string(body)) //nolint:goerr113
	}

	return nil
}

func UpdateConfig(agentName, cfgAssignmentID, cfgFile string, data []byte) error {
	trackerFile, err := getTrackerFile(agentName, cfgFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	t, err := loadTracker(trackerFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if t.AgentID == "" {
		id, err := registration.GetInstalledAgentID(agentName) //nolint:govet
		if err != nil {
			return err //nolint:wrapcheck
		}

		t.AgentID = id
	}

	t.AssignmentID = cfgAssignmentID
	t.Modified = false

	s, err := generateChecksum(cfgFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	t.S = s
	t.D = base64.StdEncoding.EncodeToString(data)

	if err := saveTracker(trackerFile, t); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func getTrackerFile(agentName, cfgFile string) (string, error) {
	baseDir, err := getBasePath(agentName)
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return filepath.Join(baseDir, filepath.Base(cfgFile)+".current"+".yaml"), nil
}

func loadTracker(file string) (*Tracker, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err //nolint:wrapcheck
		}

		return &Tracker{}, nil
	}

	var t Tracker
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &t, nil
}

func saveTracker(file string, t *Tracker) error {
	data, err := yaml.Marshal(t)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if err := os.WriteFile(file, data, 0600); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func getBasePath(agentID string) (string, error) {
	baseDir := filepath.Join(defaults.EtcPath, "configs", agentID)
	// e.g. /opt/circonus/am/etc/configs/telegraf

	if err := os.MkdirAll(baseDir, 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Error().Err(err).Str("path", baseDir).Msg("unable to make config dir to save backup")

			return "", err //nolint:wrapcheck
		}
	}

	return baseDir, nil
}

func generateChecksum(cfgFile string) (string, error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		return "", err //nolint:wrapcheck
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err //nolint:wrapcheck
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
