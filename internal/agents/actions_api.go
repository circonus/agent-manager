package agents

import (
	"encoding/json"
	"fmt"

	"github.com/circonus/agent-manager/internal/env"
	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/rs/zerolog/log"
)

type APIActions []APIAction

type APIAction struct {
	ConfigAssignmentID string         `json:"config_assignment_id"`
	Config             APIConfig      `json:"configuration"`
	Agent              APIConfigAgent `json:"agent"`
}

type APIConfig struct {
	FileID   string `json:"config_file_id"`
	Contents string `json:"config"`
}

type APIConfigAgent struct {
	ID string `json:"agent_type_id"`
}

func ParseAPIActions(data []byte) (Actions, error) {
	agents, err := inventory.LoadAgents()
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

	platform := env.GetPlatform()

	for _, apiAction := range apiActions {
		coll, ok := agents[platform][apiAction.Agent.ID]
		if !ok {
			log.Warn().Str("agent", apiAction.Agent.ID).Str("platform", platform).Msg("unknown agent for this platform")

			continue
		}

		file, ok := coll.ConfigFiles[apiAction.Config.FileID]
		if !ok {
			log.Warn().Str("id", apiAction.Config.FileID).Msg("unknown config file id")

			continue
		}

		if _, ok := actions[0].Configs[apiAction.Agent.ID]; !ok {
			actions[0].Configs = make(map[string][]Config)
		}

		cfgs := actions[0].Configs[apiAction.Agent.ID]
		cfgs = append(cfgs, Config{
			ID:       apiAction.ConfigAssignmentID,
			Path:     file,
			Contents: apiAction.Config.Contents,
		})
		actions[0].Configs[apiAction.Agent.ID] = cfgs

		foundConfigs++
	}

	if foundConfigs == 0 {
		return Actions{}, fmt.Errorf("no configs found to install")
	}

	return actions, nil
}
