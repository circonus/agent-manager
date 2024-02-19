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
	"github.com/circonus/agent-manager/internal/decommission"
	"github.com/circonus/agent-manager/internal/env"
	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/circonus/agent-manager/internal/registration"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/circonus/agent-manager/internal/server"
	"github.com/circonus/agent-manager/internal/tracker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Manager holds the main manager process.
type Manager struct {
	group       *errgroup.Group
	groupCtx    context.Context
	groupCancel context.CancelFunc
	signalCh    chan os.Signal
	logger      zerolog.Logger
	server      *server.Server
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

	if viper.GetBool(keys.Decommission) {
		if err := decommission.Start(m.groupCtx); err != nil {
			log.Fatal().Err(err).Msg("decommissioning agent manager")
		} else {
			m.logger.Info().Msg("decommission complete")
			os.Exit(0)
		}
	}

	if env.IsRunningInDocker() {
		viper.Set(keys.UseMachineID, false)
	}

	// initial registration status
	isRegistered := registration.IsRegistered()

	if viper.GetString(keys.Register) != "" {
		if isRegistered {
			log.Info().Msg("agent manager already registered, see --force-register")
		} else {
			if err := registration.Start(m.groupCtx); err != nil {
				log.Fatal().Err(err).Msg("registering agent manager")
			}
		}
	}

	// ensure manager is registered
	if !registration.IsRegistered() {
		log.Fatal().Msg("manager not registered, see instructions for registration")
	}

	if err := credentials.LoadManagerID(); err != nil {
		log.Fatal().Err(err).Msg("loading manager id")
	}

	if err := credentials.LoadJWT(); err != nil {
		log.Fatal().Err(err).Msg("loading API credentials")
	}

	if viper.GetString(keys.Register) != "" && env.IsRunningInDocker() {
		// verify that --agents and --instance-id have been provided when running in docker
		if len(viper.GetStringSlice(keys.Agents)) == 0 {
			log.Fatal().Msg("--agents required to run in container")
		}

		if viper.GetString(keys.InstanceID) == "" {
			log.Fatal().Msg("--instance-id required to run in a container")
		}
	}

	if !env.IsRunningInDocker() {
		//
		// these two are command line actions and will exit after completion
		// when not running in a docker/container.
		//
		if viper.GetString(keys.Register) != "" {
			m.logger.Info().Msg("registration complete")
			os.Exit(0)
		}
	}

	//
	// these run every time the manager starts
	//

	if err := registration.UpdateVersion(m.groupCtx); err != nil {
		m.logger.Warn().Err(err).Msg("updating manager version via API")
	}

	if err := inventory.FetchAgents(m.groupCtx); err != nil {
		log.Fatal().Err(err).Msg("fetching agents")
	}

	if err := inventory.CheckForAgents(m.groupCtx); err != nil {
		log.Fatal().Err(err).Msg("checking for installed agents")
	}

	m.logger.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("starting wait")

	actionPoller, err := agents.NewActionPoller()
	if err != nil {
		m.logger.Fatal().Err(err).Msg("unable to start action poller")
	}

	trackerPoller, err := tracker.NewPoller()
	if err != nil {
		m.logger.Fatal().Err(err).Msg("unable to start config tracker poller")
	}

	server, err := server.New()
	if err != nil {
		m.logger.Fatal().Err(err).Msg("unable to start server")
	}

	m.server = server

	m.group.Go(func() error {
		actionPoller.Start(m.groupCtx)

		return nil
	})

	m.group.Go(func() error {
		trackerPoller.Start(m.groupCtx)

		return nil
	})

	m.group.Go(func() error {
		return server.Start(m.groupCtx)
	})

	// if not running in docker, start the agent status poller
	if !env.IsRunningInDocker() {
		m.logger.Info().Msg("starting agent status poller")

		statusPoller, err := agents.NewStatusPoller()
		if err != nil {
			m.logger.Fatal().Err(err).Msg("unable to agent start status poller")
		}

		m.group.Go(func() error {
			statusPoller.Start(m.groupCtx)

			return nil
		})
	} else {
		m.logger.Info().Msg("NOT starting agent status poller -- running in docker or container")
	}

	if err := m.group.Wait(); err != nil {
		return fmt.Errorf("start manager: %w", err)
	}

	return nil
}

// Stop cleans up and shuts down the manager.
func (m *Manager) Stop() {
	m.stopSignalHandler()

	if err := m.server.Stop(m.groupCtx); err != nil {
		m.logger.Warn().Err(err).Msg("stopping server")
	}

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
