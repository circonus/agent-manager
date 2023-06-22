package collectors

import (
	"context"
	"encoding/base64"
	"os"
	"runtime"

	"github.com/rs/zerolog/log"
)

func installConfigs(ctx context.Context, a Action) {
	collectors, err := LoadCollectors()
	if err != nil {
		log.Warn().Err(err).Msg("unable to load collectors, skipping configs")
		return
	}
	for collector, configs := range a.Configs {
		for _, config := range configs {
			log.Debug().Str("path", config.Path).Str("contents", config.Contents).Msg("incoming contents")

			data, err := base64.StdEncoding.DecodeString(config.Contents)
			if err != nil {
				r := ConfigResult{
					ID: config.ID,
					ConfigData: ConfigData{
						WriteResult: err.Error(),
					},
				}
				if err = sendConfigResult(ctx, r); err != nil {
					log.Error().Err(err).Msg("config result")
				}
				continue
			}

			log.Debug().Str("path", config.Path).Str("contents", string(data)).Msg("decoded contents")

			perms := os.FileMode(0640)

			f, err := os.Stat(config.Path)
			if err == nil {
				perms = f.Mode().Perm()
			}

			if err := os.WriteFile(config.Path, data, perms); err != nil {
				r := ConfigResult{
					ID:     config.ID,
					Status: STATUS_ERROR,
					Info:   err.Error(),
					ConfigData: ConfigData{
						WriteResult: err.Error(),
					},
				}
				if err := sendConfigResult(ctx, r); err != nil {
					log.Error().Err(err).Msg("config result")
				}
				continue
			}

			r := ConfigResult{
				ID:     config.ID,
				Status: STATUS_ACTIVE,
				ConfigData: ConfigData{
					WriteResult: "OK",
				},
			}
			if err := sendConfigResult(ctx, r); err != nil {
				log.Error().Err(err).Msg("config result")
			}
		}
		c, ok := collectors[runtime.GOOS][collector]
		if !ok {
			log.Warn().Str("platform", runtime.GOOS).Str("collector", collector).Msg("unable to find collector definition for reload, skipping")
			continue
		}
		if c.Reload == "" {
			continue
		}
		if c.Reload == RESTART {
			output, code, err := execute(ctx, c.Restart)
			if err != nil {
				log.Warn().Err(err).Str("output", string(output)).Int("exit_code", code).Str("cmd", c.Restart).Msg("restart failed")
			}
		}
	}
}
