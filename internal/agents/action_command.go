// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"context"
	"encoding/base64"

	"github.com/circonus/agent-manager/internal/inventory"
	"github.com/circonus/agent-manager/internal/platform"
	"github.com/rs/zerolog/log"
)

const (
	START     = "start"
	STOP      = "stop"
	RESTART   = "restart"
	STATUS    = "status"
	RELOAD    = "reload"
	INVENTORY = "inventory"
	VERSION   = "version"
)

func runCommands(ctx context.Context, action Action) error {
	agents, err := inventory.LoadAgents()
	if err != nil {
		return err // nolint:wrapcheck
	}

	platform := platform.Get()

	for _, command := range action.Commands {
		switch command.Command {
		case INVENTORY:
			if err := inventory.FetchAgents(ctx); err != nil {
				log.Error().Err(err).Msg("refreshing agent list")
			} else if err := inventory.CheckForAgents(ctx); err != nil {
				log.Error().Err(err).Msg("checking for installed agents")
			}
		case START:
			a, ok := agents[platform][command.Agent]
			if ok {
				runCommand(ctx, a.Start, command.ID)
			}
		case STOP:
			a, ok := agents[platform][command.Agent]
			if ok {
				runCommand(ctx, a.Stop, command.ID)
			}
		case RESTART:
			a, ok := agents[platform][command.Agent]
			if ok {
				runCommand(ctx, a.Restart, command.ID)
			}
		case RELOAD:
			a, ok := agents[platform][command.Agent]
			if ok {
				cmdReload(ctx, a, command)
			}
		case STATUS:
			a, ok := agents[platform][command.Agent]
			if ok {
				runCommand(ctx, a.Status, command.ID)
			}
		case VERSION:
			a, ok := agents[platform][command.Agent]
			if ok {
				runCommand(ctx, a.Version, command.ID)
			}
		}
	}

	return nil
}

func runCommand(ctx context.Context, cmd, id string) {
	output, code, err := execute(ctx, cmd)
	if err != nil {
		log.Warn().Err(err).Str("output", string(output)).Int("exit_code", code).Str("cmd", cmd).Msg("command failed")
	}

	if id != "" {
		result := CommandResult{
			ID: id,
			CommandData: CommandData{
				ExitCode: code,
			},
		}

		if err != nil {
			result.CommandData.Error = err.Error()
		}

		if len(output) > 0 {
			result.CommandData.Output = base64.StdEncoding.EncodeToString(output)
		}

		if err = sendCommandResult(ctx, result); err != nil {
			log.Error().Err(err).Msg("command result")
		}
	}
}
