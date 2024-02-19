package agents

import (
	"context"
	"encoding/base64"

	"github.com/circonus/agent-manager/internal/env"
	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/circonus/agent-manager/internal/server"
	"github.com/circonus/agent-manager/internal/tracker"
	"github.com/rs/zerolog/log"
)

func installConfigs(ctx context.Context, action Action) {
	agents, err := inventory.LoadAgents()
	if err != nil {
		log.Warn().Err(err).Msg("unable to load agents, skipping configs")

		return
	}

	platform := env.GetPlatform()

	for agentID, configs := range action.Configs {
		for _, config := range configs {
			log.Debug().Str("path", config.Path).Str("contents", config.Contents).Msg("incoming contents")

			data, err := base64.StdEncoding.DecodeString(config.Contents)
			if err != nil {
				result := ConfigResult{
					ID: config.ID,
					ConfigData: ConfigData{
						WriteResult: err.Error(),
					},
				}

				if err = sendConfigResult(ctx, result); err != nil {
					log.Error().Err(err).Msg("config result")
				}

				continue
			}

			log.Debug().Str("path", config.Path).Str("contents", string(data)).Msg("decoded contents")

			if err := writeConfig(config.Path, data); err != nil {
				result := ConfigResult{
					ID:     config.ID,
					Status: STATUS_ERROR,
					Info:   err.Error(),
					ConfigData: ConfigData{
						WriteResult: err.Error(),
					},
				}

				if err := sendConfigResult(ctx, result); err != nil {
					log.Error().Err(err).Msg("config result")
				}

				continue
			}

			result := ConfigResult{
				ID:     config.ID,
				Status: STATUS_ACTIVE,
				ConfigData: ConfigData{
					WriteResult: "OK",
				},
			}

			if err := sendConfigResult(ctx, result); err != nil {
				log.Error().Err(err).Msg("config result")
			}

			// save config hash as current.
			if err := tracker.UpdateConfig(agentID, config.ID, config.Path, data); err != nil {
				log.Error().Err(err).Msg("updating config tracking data")
			}
		}

		if env.IsRunningInDocker() {
			server.AddConfigUpdate(agentID)
		} else {
			agent, ok := agents[platform][agentID]
			if !ok {
				log.Warn().Str("platform", platform).Str("agent", agentID).
					Msg("unable to find agent definition for reload, skipping")

				continue
			}

			cmdReload(ctx, agent, Command{})
		}
	}
}
