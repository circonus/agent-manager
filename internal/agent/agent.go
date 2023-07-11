// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/circonus/agent-manager/internal/collectors"
	"github.com/circonus/agent-manager/internal/config"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/credentials"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Agent holds the main agent process.
type Agent struct {
	group       *errgroup.Group
	groupCtx    context.Context //nolint:containedctx
	groupCancel context.CancelFunc
	signalCh    chan os.Signal
	logger      zerolog.Logger
}

// New returns a new agent instance.
func New() (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	var err error

	agent := Agent{
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

	agent.signalNotifySetup()

	return &agent, nil
}

// Start is the main agent entry point.
func (a *Agent) Start() error {
	a.group.Go(a.handleSignals)

	log.Info().Str("name", release.NAME).Str("version", release.VERSION).Msg("starting")

	if viper.GetString(keys.Register) != "" {
		if err := credentials.SaveRegistrationToken([]byte(viper.GetString(keys.Register))); err != nil {
			log.Fatal().Err(err).Msg("saving registration token")
		}

		if err := registration.Start(a.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("registering agent")
		}
	}

	if err := credentials.LoadManagerID(); err != nil {
		log.Fatal().Err(err).Msg("loading manager id")
	}

	if err := credentials.LoadRegistrationToken(); err != nil {
		log.Fatal().Err(err).Msg("loading registration token")
	}

	if err := credentials.LoadJWT(); err != nil {
		log.Fatal().Err(err).Msg("loading API credentials")
	}

	if viper.GetString(keys.Register) != "" || viper.GetBool(keys.Inventory) {
		if err := collectors.FetchCollectors(a.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("fetching collectors")
		}

		if err := collectors.CheckForCollectors(a.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("checking for installed collectors")
		}
	}

	a.logger.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("starting wait")

	poller, err := collectors.NewPoller()
	if err != nil {
		a.logger.Fatal().Err(err).Msg("unable to start poller")
	}

	a.group.Go(func() error {
		poller.Start(a.groupCtx)

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
