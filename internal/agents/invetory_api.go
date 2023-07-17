// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type APIAgents []APIAgent

type APIAgent struct {
	Platforms []Platform `json:"platforms"`
}

type Platform struct {
	ID          string       `json:"platform_id"   yaml:"platform_id"`
	AgentTypeID string       `json:"agent_type_id" yaml:"agent_type_id"`
	Executable  string       `json:"executable"    yaml:"executable"`
	Start       string       `json:"start"         yaml:"start"`
	Stop        string       `json:"stop"          yaml:"stop"`
	Restart     string       `json:"restart"       yaml:"restart"`
	Reload      string       `json:"reload"        yaml:"reload"`
	Status      string       `json:"status"        yaml:"status"`
	Version     string       `json:"version"       yaml:"version"`
	ConfigFiles []ConfigFile `json:"config_files"  yaml:"config_files"`
}

type ConfigFile struct {
	ConfigFileID string `json:"config_file_id" yaml:"config_file_id"`
	Path         string `json:"path"           yaml:"path"`
}

func ParseAPIAgents(data []byte) (Agents, error) {
	var aa APIAgents
	if err := json.Unmarshal(data, &aa); err != nil {
		log.Error().Err(err).Msg("parsing api agents json")
	}

	agents := make(map[string]map[string]Agent)

	for _, agent := range aa {
		for _, platform := range agent.Platforms {
			col := Agent{
				Binary:      platform.Executable,
				Start:       platform.Start,
				Stop:        platform.Stop,
				Restart:     platform.Restart,
				Reload:      platform.Reload,
				Version:     platform.Version,
				ConfigFiles: make(map[string]string, len(platform.ConfigFiles)),
			}

			for _, f := range platform.ConfigFiles {
				col.ConfigFiles[f.ConfigFileID] = f.Path
			}

			if _, ok := agents[platform.ID]; !ok {
				agents[platform.ID] = make(map[string]Agent)
			}

			agents[platform.ID][platform.AgentTypeID] = col
		}
	}

	return agents, nil
}
