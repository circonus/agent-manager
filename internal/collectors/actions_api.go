// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
)

type APIActions []APIAction

type APIAction struct {
	ConfigAssignmentID string             `json:"collector_config_assignment_id"`
	Config             APIConfig          `json:"collector_config"`
	Collector          APIConfigCollector `json:"collector"`
}

type APIConfig struct {
	FileID   string `json:"collector_config_file_id"`
	Contents string `json:"config"`
}

type APIConfigCollector struct {
	ID string `json:"collector_type_id"`
}

func ParseAPIActions(data []byte) (Actions, error) {
	collectors, err := LoadCollectors()
	if err != nil {
		return nil, err
	}

	var apiActions APIActions
	if err := json.Unmarshal(data, &apiActions); err != nil {
		return nil, fmt.Errorf("parsing api actions: %w", err)
	}

	actions := []Action{
		{
			Type:    CONFIG,
			Configs: make(map[string][]Config),
		},
	}

	foundConfigs := 0

	for _, apiAction := range apiActions {
		coll, ok := collectors[runtime.GOOS][apiAction.Collector.ID]
		if !ok {
			log.Warn().Str("collector", apiAction.Collector.ID).Str("platform", runtime.GOOS).Msg("unknown collector for this platform")

			continue
		}

		file, ok := coll.ConfigFiles[apiAction.Config.FileID]
		if !ok {
			log.Warn().Str("id", apiAction.Config.FileID).Msg("unknown config file id")

			continue
		}

		if _, ok := actions[0].Configs[apiAction.Collector.ID]; !ok {
			actions[0].Configs = make(map[string][]Config)
		}

		cfgs := actions[0].Configs[apiAction.Collector.ID]
		cfgs = append(cfgs, Config{
			ID:       apiAction.ConfigAssignmentID,
			Path:     file,
			Contents: apiAction.Config.Contents,
		})
		actions[0].Configs[apiAction.Collector.ID] = cfgs

		foundConfigs++
	}

	if foundConfigs == 0 {
		return Actions{}, fmt.Errorf("no configs found to install") //nolint:goerr113
	}

	return actions, nil
}
