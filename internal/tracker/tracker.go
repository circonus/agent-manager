// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package tracker

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Tracker struct {
	AgentID string `json:"agent_id" yaml:"agent_id"`
	S       string `json:"s"        yaml:"s"`
	D       string `json:"d"        yaml:"d"`
}

func VerifyConfig(agentName, cfgID, cfgFile string) error {
	trackerFile, err := getTrackerFile(agentName, cfgFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	t, err := loadTracker(trackerFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if t.AgentID == "" || t.S == "" {
		return fmt.Errorf("no current tracking information available") //nolint:goerr113
	}

	s, err := generateChecksum(cfgFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if s != t.S {
		return fmt.Errorf("tracking signature different: c[%s] v[%s] id:%s", s, t.S, cfgID) //nolint:goerr113
	}

	return nil
}

func UpdateConfig(agentName, cfgFile string, data []byte) error {
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
		if errors.Is(err, os.ErrNotExist) {
			return &Tracker{}, nil
		}

		return nil, err //nolint:wrapcheck
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

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
