// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"context"
	"encoding/base64"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
)

func installConfigs(ctx context.Context, action Action) {
	collectors, err := LoadCollectors()
	if err != nil {
		log.Warn().Err(err).Msg("unable to load collectors, skipping configs")

		return
	}

	for collector, configs := range action.Configs {
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
					Status: STATUS_ERROR, //nolint:nosnakecase
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
				Status: STATUS_ACTIVE, //nolint:nosnakecase
				ConfigData: ConfigData{
					WriteResult: "OK",
				},
			}

			if err := sendConfigResult(ctx, result); err != nil {
				log.Error().Err(err).Msg("config result")
			}
		}

		coll, ok := collectors[runtime.GOOS][collector]
		if !ok {
			log.Warn().Str("platform", runtime.GOOS).Str("collector", collector).Msg("unable to find collector definition for reload, skipping")

			continue
		}

		switch {
		case coll.Reload == "":
			continue
		case strings.ToLower(coll.Reload) == RESTART:
			output, code, err := execute(ctx, coll.Restart)
			if err != nil {
				log.Warn().Err(err).Str("output", string(output)).Int("exit_code", code).Str("cmd", coll.Restart).Msg("restart failed")
			}
		//FUTURE: add case(s) for other options e.g. hitting an endpoint to trigger a reload
		default:
			output, code, err := execute(ctx, coll.Reload)
			if err != nil {
				log.Warn().Err(err).Str("output", string(output)).Int("exit_code", code).Str("cmd", coll.Reload).Msg("reload failed")
			}
		}
	}
}
