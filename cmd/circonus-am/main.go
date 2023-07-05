// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build go1.20
// +build go1.20

package main

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime"

	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string //nolint:gochecknoglobals

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zlog := zerolog.New(zerolog.SyncWriter(os.Stderr)).With().Timestamp().Logger()
	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)

	cobra.OnInitialize(initConfig)

	cmd := initCmd()
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("starting agent")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(defaults.EtcPath)
		viper.AddConfigPath(".")
		viper.SetConfigName(release.NAME)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		f := viper.ConfigFileUsed()
		if f != "" {
			log.Fatal().Err(err).Str("config_file", f).Msg("unable to load config file")
		}
	}
}

// initApp initializes the application components.
func initApp(_ *cobra.Command, _ []string) error {
	return initLogging()
}

// initLogging initializes zerolog.
func initLogging() error {
	//
	// Enable formatted output
	//
	if viper.GetBool(keys.LogPretty) {
		if runtime.GOOS != "windows" {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		} else {
			log.Warn().Msg("log-pretty not applicable on this platform")
		}
	}

	//
	// Enable debug logging if requested
	//
	if viper.GetBool(keys.Debug) {
		viper.Set(keys.LogLevel, "debug")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		return nil
	}

	//
	// otherwise, set custom level if specified
	//
	if viper.IsSet(keys.LogLevel) {
		level := viper.GetString(keys.LogLevel)

		switch level {
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "disabled":
			zerolog.SetGlobalLevel(zerolog.Disabled)
		default:
			return fmt.Errorf("unknown log level (%s)", level) //nolint:goerr113
		}
	}

	return nil
}
