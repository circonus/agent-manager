// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"fmt"
	"runtime"

	"github.com/circonus/collector-management-agent/internal/agent"
	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/circonus/collector-management-agent/internal/release"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   release.NAME,
		Short: "A brief description of this agent",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using this agent.`,
		PersistentPreRunE: initApp,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(keys.ShowVersion) {
				fmt.Printf("%s v%s - commit: %s, date: %s, tag: %s, built with: %s\n", release.NAME, release.VERSION, release.COMMIT, release.DATE, release.TAG, runtime.Version())
				return
			}

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
