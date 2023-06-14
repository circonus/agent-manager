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
	for {
		t := time.NewTimer(p.interval)
		select {
		case <-ctx.Done():
			if !t.Stop() {
				<-t.C
			}
			return
		case <-t.C:
			if err := getActions(ctx); err != nil {
				log.Error().Err(err).Msg("getting actions")
			}
		}
	}
}
