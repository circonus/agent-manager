package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/circonus/go-agent-template/internal/config"
	"github.com/circonus/go-agent-template/internal/config/keys"
	"github.com/circonus/go-agent-template/internal/release"
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

	log.Info().Str("example_arg", viper.GetString(keys.ExampleArg)).Msg("example argument")

	a.logger.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("starting wait")

	// if err := a.group.Wait(); err != nil {
	// 	return fmt.Errorf("start agent: %w", err)
	// }

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
