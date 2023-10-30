// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/env"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

/*
  1. provide generic /health endpoint.
  2. provide docker containers a health check endpoint /config/<agent>
     used to trigger container process reloads on configuration file changes.
*/

type Server struct {
	srv             *http.Server
	idleConnsClosed chan struct{}
}

func New() (*Server, error) {
	readTimeout, err := time.ParseDuration(viper.GetString(keys.ServerReadTimeout))
	if err != nil {
		return nil, fmt.Errorf("parsing read timeout (%s): %w", viper.GetString(keys.ServerReadTimeout), err)
	}

	writeTimeout, err := time.ParseDuration(viper.GetString(keys.ServerWriteTimeout))
	if err != nil {
		return nil, fmt.Errorf("parsing write timeout (%s): %w", viper.GetString(keys.ServerWriteTimeout), err)
	}

	idleTimeout, err := time.ParseDuration(viper.GetString(keys.ServerIdleTimeout))
	if err != nil {
		return nil, fmt.Errorf("parsing idle timeout (%s): %w", viper.GetString(keys.ServerIdleTimeout), err)
	}

	readHeaderTimeout, err := time.ParseDuration(viper.GetString(keys.ServerReadHeaderTimeout))
	if err != nil {
		return nil, fmt.Errorf("parsing read header timeout (%s): %w", viper.GetString(keys.ServerReadHeaderTimeout), err)
	}

	handlerTimeout, err := time.ParseDuration(viper.GetString(keys.ServerHandlerTimeout))
	if err != nil {
		return nil, fmt.Errorf("parsing handler timeout (%s): %w", viper.GetString(keys.ServerHandlerTimeout), err)
	}

	mux := http.NewServeMux()

	mux.Handle("/health", reqLogger(http.TimeoutHandler(
		healthHandler{}, handlerTimeout, "health handler timeout")))

	if env.IsRunningInDocker() {
		// e.g. Docker, when a config has changed, /config will return a 409 (conflict),
		//      indicating that the agent should reload its config(s)
		//
		//   HEALTHCHECK --interval=90s --timeout=3s \
		//     CMD curl --silent --fail "http://<cam-container-ip>:43285/config/<agent_type>" || exit 1
		//     CMD wget --quiet "http://<cam-container-ip>:43285/config/<agent_type>" || exit 1
		mux.Handle("/config", reqLogger(http.TimeoutHandler(
			configHandler{}, handlerTimeout, "config handler timeout")))
	}

	return &Server{
		srv: &http.Server{
			Addr:              viper.GetString(keys.ServerAddress),
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			Handler:           mux,
		},
		idleConnsClosed: make(chan struct{}),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	if done(ctx) {
		return ctx.Err() //nolint:wrapcheck
	}

	if viper.GetBool(keys.ServerTLSEnable) && //nolint:nestif
		viper.GetString(keys.ServerTLSCertFile) != "" &&
		viper.GetString(keys.ServerTLSKeyFile) != "" {
		log.Info().Str("listen", s.srv.Addr).Msg("starting TLS server")

		if err := s.srv.ListenAndServeTLS(
			viper.GetString(keys.ServerTLSCertFile),
			viper.GetString(keys.ServerTLSKeyFile)); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("listen and serve tls")
			}
		}
	} else {
		log.Info().Str("listen", s.srv.Addr).Msg("starting server")
		if err := s.srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("listen and serve")
			}
		}
	}

	<-s.idleConnsClosed

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info().Msg("shutting down server")

	toctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(toctx); err != nil {
		log.Error().Err(err).Msg("server shutdown")
	}

	close(s.idleConnsClosed)

	// if no error, check the ctx and return that error
	if done(ctx) {
		return ctx.Err() //nolint:wrapcheck
	}

	return nil
}

func done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func reqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("url", r.URL.String()).Msg("incoming request")

		next.ServeHTTP(w, r)
	})
}
