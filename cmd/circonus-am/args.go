package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func bindFlagError(flag string, err error) {
	if err != nil {
		log.Fatal().Err(err).Str("flag", flag).Msg("binding flag")
	}
}

func bindEnvError(envVar string, err error) {
	if err != nil {
		log.Fatal().Err(err).Str("var", envVar).Msg("binding env var")
	}
}

func envDescription(desc, env string) string {
	if env == "" {
		return desc
	}

	return fmt.Sprintf("[ENV: %s] %s", env, desc)
}

func initArgs(cmd *cobra.Command) {
	initGeneralArgs(cmd)
	initAppArgs(cmd)
}
