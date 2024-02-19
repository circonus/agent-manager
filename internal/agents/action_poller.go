package agents

import (
	"context"
	"fmt"
	"time"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// manages polling for actions

type ActionPoller struct {
	interval time.Duration
}

func NewActionPoller() (*ActionPoller, error) {
	pi := viper.GetString(keys.ActionPollingInterval)

	i, err := time.ParseDuration(pi)
	if err != nil {
		return nil, fmt.Errorf("parsing polling interval: %w", err)
	}

	return &ActionPoller{interval: i}, nil
}

func (p *ActionPoller) Start(ctx context.Context) {
	log.Info().Str("interval", p.interval.String()).Msg("starting action poller")

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
