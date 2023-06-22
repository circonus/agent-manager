// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

import (
	"context"
	"fmt"
	"time"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// manages polling for actions

type Poller struct {
	interval time.Duration
}

func NewPoller() (*Poller, error) {
	pi := viper.GetString(keys.PollingInterval)

	i, err := time.ParseDuration(pi)
	if err != nil {
		return nil, fmt.Errorf("parsing polling interval: %w", err)
	}

	return &Poller{interval: i}, nil
}

func (p *Poller) Start(ctx context.Context) {
	log.Info().Str("interval", p.interval.String()).Msg("starting poller")
	for {
		t := time.NewTimer(p.interval)
		select {
		case <-ctx.Done():
			if !t.Stop() {
				<-t.C
			}
			return
		case <-t.C:
			log.Debug().Msg("checking for new actions")
			if err := getActions(ctx); err != nil {
				log.Error().Err(err).Msg("getting actions")
			}
		}
	}
}
