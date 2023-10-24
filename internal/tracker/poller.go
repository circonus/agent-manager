// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Poller struct {
	interval time.Duration
}

func NewPoller() (*Poller, error) {
	pi := viper.GetString(keys.TrackerPollingInterval)

	i, err := time.ParseDuration(pi)
	if err != nil {
		return nil, fmt.Errorf("parsing tracker polling interval: %w", err)
	}

	return &Poller{interval: i}, nil
}

func (p *Poller) Start(ctx context.Context) {
	log.Info().Str("interval", p.interval.String()).Msg("starting config tracker")

	for {
		t := time.NewTimer(p.interval)
		select {
		case <-ctx.Done():
			if !t.Stop() {
				<-t.C
			}

			return
		case <-t.C:
			log.Debug().Msg("tracking installed configs")

			agents, err := registration.LoadInstalledAgents()
			if err != nil {
				log.Error().Err(err).Msg("loading installed agents, run --inventory again to generate")

				continue
			}

			for _, a := range agents {
				agent, err := inventory.GetAgent(a.AgentTypeID)
				if err != nil {
					log.Error().Err(err).Str("agent_type", a.AgentTypeID).Msg("getting agent")

					continue
				}

				for cfgID, path := range agent.ConfigFiles {
					if err := VerifyConfig(ctx, a.AgentTypeID, path); err != nil {
						log.Error().Err(err).
							Str("agent", a.AgentTypeID).
							Str("id", cfgID).
							Str("file", path).
							Msg("config tracking issue")

						continue
					}
				}
			}
		}
	}
}
