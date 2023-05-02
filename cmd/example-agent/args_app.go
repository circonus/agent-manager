// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"github.com/circonus/go-agent-template/internal/config/defaults"
	"github.com/circonus/go-agent-template/internal/config/keys"
	"github.com/circonus/go-agent-template/internal/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initAppArgs adds application specific args to the cobra command.
func initAppArgs(cmd *cobra.Command) {
	{
		const (
			key          = keys.ExampleArg
			longOpt      = "example-arg"
			envVar       = release.ENVPREFIX + "_EXAMPLE_ARG"
			description  = "Example Argument"
			defaultValue = defaults.ExampleArg
		)

		cmd.PersistentFlags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.PersistentFlags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

}
