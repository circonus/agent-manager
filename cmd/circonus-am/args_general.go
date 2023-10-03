// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initGeneralArgs adds general args to the cobra command.
func initGeneralArgs(cmd *cobra.Command) {
	{
		var (
			longOpt     = "config"
			shortOpt    = "c"
			description = "config file (default: " + defaults.ConfigFile + "|.json|.toml)"
		)
		cmd.Flags().StringVarP(&cfgFile, longOpt, shortOpt, "", description)
	}

	{
		const (
			key          = keys.ShowVersion
			longOpt      = "version"
			shortOpt     = "V"
			defaultValue = false
			description  = "Show version and exit"
		)
		cmd.Flags().BoolP(longOpt, shortOpt, defaultValue, description)
		if err := viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)); err != nil {
			bindFlagError(longOpt, err)
		}
	}

	{
		const (
			key          = keys.Debug
			longOpt      = "debug"
			shortOpt     = "d"
			envVar       = release.ENVPREFIX + "_DEBUG"
			description  = "Enable debug messages"
			defaultValue = defaults.Debug
		)

		cmd.Flags().BoolP(longOpt, shortOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.LogLevel
			longOpt      = "log-level"
			envVar       = release.ENVPREFIX + "_LOG_LEVEL"
			description  = "Log level [(panic|fatal|error|warn|info|debug|disabled)]"
			defaultValue = defaults.LogLevel
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.LogPretty
			longOpt      = "log-pretty"
			description  = "Output formatted/colored log lines [ignored on windows]"
			defaultValue = defaults.LogPretty
		)

		cmd.Flags().Bool(longOpt, defaultValue, description)
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		viper.SetDefault(key, defaultValue)
	}
}
