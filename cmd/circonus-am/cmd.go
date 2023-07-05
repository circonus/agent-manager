// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"fmt"
	"runtime"

	"github.com/circonus/agent-manager/internal/agent"
	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               release.NAME,
		Short:             "Circonus Management Agent",
		Long:              `An agent to manage local collectors (metrics, logs, etc.)`,
		PersistentPreRunE: initApp,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(keys.ShowVersion) {
				fmt.Printf("%s v%s - commit: %s, date: %s, tag: %s, built with: %s\n", release.NAME, release.VERSION, release.COMMIT, release.DATE, release.TAG, runtime.Version())

				return
			}

			// set internal viper items
			viper.Set(keys.InventoryFile, defaults.InventoryFile)
			viper.Set(keys.JwtTokenFile, defaults.JwtTokenFile)
			viper.Set(keys.RegTokenFile, defaults.RegTokenFile)
			viper.Set(keys.AgentIDFile, defaults.AgentIDFile)

			a, err := agent.New()
			if err != nil {
				log.Fatal().Err(err).Msg("initializing")
			}

			if err := a.Start(); err != nil {
				log.Fatal().Err(err).Msg("starting agent")
			}
		},
	}

	initArgs(cmd)

	return cmd
}
