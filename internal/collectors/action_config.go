package collectors

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/rs/zerolog/log"
)

func installConfigs(ctx context.Context, a Action) {
	for _, config := range a.Configs {
		data, err := base64.StdEncoding.DecodeString(config.Contents)
		if err != nil {
			r := Result{
				ActionID: a.ID,
				ConfigResult: ConfigResult{
					ID: config.ID,
					ConfigData: ConfigData{
						WriteResult: err.Error(),
					},
				},
			}
			if err = sendActionResult(ctx, r); err != nil {
				log.Error().Err(err).Msg("config result")
			}
			continue
		}

		perms := os.FileMode(0640)

		f, err := os.Stat(config.Path)
		if err == nil {
			perms = f.Mode().Perm()
		}

		if err := os.WriteFile(config.Path, data, perms); err != nil {
			r := Result{
				ActionID: a.ID,
				ConfigResult: ConfigResult{
					ID: config.ID,
					ConfigData: ConfigData{
						WriteResult: err.Error(),
					},
				},
			}
			if err := sendActionResult(ctx, r); err != nil {
				log.Error().Err(err).Msg("config result")
			}
			continue
		}

		r := Result{
			ActionID: a.ID,
			ConfigResult: ConfigResult{
				ID: config.ID,
				ConfigData: ConfigData{
					WriteResult: "OK",
				},
			},
		}
		if err := sendActionResult(ctx, r); err != nil {
			log.Error().Err(err).Msg("config result")
		}
	}
}
