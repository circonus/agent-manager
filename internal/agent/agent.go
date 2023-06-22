package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/circonus/collector-management-agent/internal/collectors"
	"github.com/circonus/collector-management-agent/internal/config"
	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/circonus/collector-management-agent/internal/credentials"
	"github.com/circonus/collector-management-agent/internal/registration"
	"github.com/circonus/collector-management-agent/internal/release"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Agent holds the main agent process.
type Agent struct {
	group       *errgroup.Group
	groupCtx    context.Context
	groupCancel context.CancelFunc
	signalCh    chan os.Signal
	logger      zerolog.Logger
}

// New returns a new agent instance.
func New() (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	var err error
	a := Agent{
		group:       g,
		groupCtx:    gctx,
		groupCancel: cancel,
		signalCh:    make(chan os.Signal, 10),
		logger:      log.With().Str("pkg", "agent").Logger(),
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}

	a.signalNotifySetup()

	return &a, nil
}

// Start is the main agent entry point.
func (a *Agent) Start() error {
	a.group.Go(a.handleSignals)

	log.Info().Str("name", release.NAME).Str("version", release.VERSION).Msg("starting")

	if viper.GetString(keys.Register) != "" {
		credentials.SaveRegistrationToken([]byte(viper.GetString(keys.Register)))
		if err := registration.Start(a.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("registering agent")
		}
	}

	if err := credentials.LoadAgentID(); err != nil {
		log.Fatal().Err(err).Msg("loading agent id")
	}

	if err := credentials.LoadJWT(); err != nil {
		log.Fatal().Err(err).Msg("loading API credentials")
	}

	a.logger.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("starting wait")

	p, err := collectors.NewPoller()
	if err != nil {
		a.logger.Fatal().Err(err).Msg("unable to start poller")
	}

	a.group.Go(func() error {
		p.Start(a.groupCtx)
		return nil
	})

	a.group.Go(func() error {
		registration.ReRegister(a.groupCtx)
		return nil
	})

	if err := a.group.Wait(); err != nil {
		return fmt.Errorf("start agent: %w", err)
	}

	return nil
}

// Stop cleans up and shuts down the Agent.
func (a *Agent) Stop() {
	a.stopSignalHandler()
	a.groupCancel()

	a.logger.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Stopped")
}

// stopSignalHandler disables the signal handler.
func (a *Agent) stopSignalHandler() {
	signal.Stop(a.signalCh)
	signal.Reset() // so a second ctrl-c will force immediate stop
}
