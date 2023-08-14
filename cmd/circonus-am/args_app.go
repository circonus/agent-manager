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

// initAppArgs adds application specific args to the cobra command.
func initAppArgs(cmd *cobra.Command) {
	{
		const (
			key          = keys.Register
			longOpt      = "register"
			envVar       = release.ENVPREFIX + "_REGISTER"
			description  = "Registration token"
			defaultValue = ""
		)

		cmd.PersistentFlags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.PersistentFlags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.Inventory
			longOpt      = "inventory"
			envVar       = release.ENVPREFIX + "_INVENTORY"
			description  = "Inventory installed agents"
			defaultValue = false
		)

		cmd.PersistentFlags().Bool(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.PersistentFlags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.APIURL
			longOpt      = "apiurl"
			envVar       = release.ENVPREFIX + "_API_URL"
			description  = "Circonus API URL"
			defaultValue = defaults.APIURL
		)

		cmd.PersistentFlags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.PersistentFlags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.PollingInterval
			longOpt      = "poll-interval"
			envVar       = release.ENVPREFIX + "_POLL_INTERVAL"
			description  = "Polling interval for actions"
			defaultValue = defaults.PollingInterval
		)

		cmd.PersistentFlags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.PersistentFlags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = keys.AWSEC2Tags
			longOpt     = "aws-ec2-tags"
			envVar      = release.ENVPREFIX + "_AWS_EC2_TAGS"
			description = "AWS EC2 tags for registration meta data"
		)
		defaultValue := defaults.AWSEC2Tags

		cmd.PersistentFlags().StringArray(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.PersistentFlags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}
}
