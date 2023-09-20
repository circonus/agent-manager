// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package manager

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/circonus/agent-manager/internal/agents"
	"github.com/circonus/agent-manager/internal/config"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/credentials"
	"github.com/circonus/agent-manager/internal/decomission"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Manager holds the main manager process.
type Manager struct {
	group       *errgroup.Group
	groupCtx    context.Context //nolint:containedctx
	groupCancel context.CancelFunc
	signalCh    chan os.Signal
	logger      zerolog.Logger
}

// New returns a new manager instance.
func New() (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	var err error

	manager := Manager{
		group:       g,
		groupCtx:    gctx,
		groupCancel: cancel,
		signalCh:    make(chan os.Signal, 10),
		logger:      log.With().Str("pkg", "manager").Logger(),
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}

	manager.signalNotifySetup()

	return &manager, nil
}

// Start is the main entry point.
func (m *Manager) Start() error {
	m.group.Go(m.handleSignals)

	log.Info().Str("name", release.NAME).Str("version", release.VERSION).Msg("starting")

	if viper.GetBool(keys.Decomission) {
		if err := decomission.Start(m.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("decommissioning agent manager")
		} else {
			m.logger.Info().Msg("decomission complete")
			os.Exit(0)
		}
	}

	if viper.GetString(keys.Register) != "" {
		if err := registration.Start(m.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("registering agent manager")
		}
	}

	if !credentials.DoesFileExist(viper.GetString(keys.JwtTokenFile)) ||
		!credentials.DoesFileExist(viper.GetString(keys.ManagerIDFile)) {
		log.Fatal().Msg("manager not registered, see instructions for registeration")
	}

	if err := credentials.LoadManagerID(); err != nil {
		log.Fatal().Err(err).Msg("loading manager id")
	}

	if err := credentials.LoadJWT(); err != nil {
		log.Fatal().Err(err).Msg("loading API credentials")
	}

	if viper.GetString(keys.Register) != "" || viper.GetBool(keys.Inventory) {
		if err := agents.FetchAgents(m.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("fetching agents")
		}

		if err := agents.CheckForAgents(m.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("checking for installed agents")
		}
	}

	//
	// these two are command line actions and will exit after completion.
	//

	if viper.GetString(keys.Register) != "" {
		m.logger.Info().Msg("registration complete")
		os.Exit(0)
	}

	if viper.GetBool(keys.Inventory) {
		m.logger.Info().Msg("invetory complete")
		os.Exit(0)
	}

	m.logger.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("starting wait")

	poller, err := agents.NewPoller()
	if err != nil {
		m.logger.Fatal().Err(err).Msg("unable to start poller")
	}

	m.group.Go(func() error {
		poller.Start(m.groupCtx)

		return nil
	})

	if err := m.group.Wait(); err != nil {
		return fmt.Errorf("start manager: %w", err)
	}

	return nil
}

// Stop cleans up and shuts down the manager.
func (m *Manager) Stop() {
	m.stopSignalHandler()
	m.groupCancel()

	m.logger.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Stopped")
}

// stopSignalHandler disables the signal handler.
func (m *Manager) stopSignalHandler() {
	signal.Stop(m.signalCh)
	signal.Reset() // so a second ctrl-c will force immediate stop
}
