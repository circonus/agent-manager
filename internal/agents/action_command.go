// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"context"
	"encoding/base64"
	"runtime"

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
	agents, err := LoadAgents()
	if err != nil {
		return err
	}

	for _, command := range action.Commands {
		switch command.Command {
		case INVENTORY:
			if err := FetchAgents(ctx); err != nil {
				log.Error().Err(err).Msg("refreshing agents")
			} else if err := CheckForAgents(ctx); err != nil {
				log.Error().Err(err).Msg("checking for agents")
			}
		case START:
			a, ok := agents[runtime.GOOS][command.Agent]
			if ok {
				cmdStart(ctx, a, command)
			}
		case STOP:
			a, ok := agents[runtime.GOOS][command.Agent]
			if ok {
				cmdStop(ctx, a, command)
			}
		case RESTART:
			a, ok := agents[runtime.GOOS][command.Agent]
			if ok {
				cmdRestart(ctx, a, command)
			}
		case RELOAD:
			// this needs to be handled differently as reload may be:
			// a command or some type of endpoint
			//
			// a, ok := agents[runtime.GOOS][command.Agent]
			// if ok {
			// args := strings.Split(a.Reload, " ")
			// cmd := exec.Command(args[0], args[1:]...)
			// }
		case STATUS:
			a, ok := agents[runtime.GOOS][command.Agent]
			if ok {
				cmdStatus(ctx, a, command)
			}
		case VERSION:
			a, ok := agents[runtime.GOOS][command.Agent]
			if ok {
				cmdVersion(ctx, a, command)
			}
		}
	}

	return nil
}

func cmdStart(ctx context.Context, agent Agent, command Command) {
	output, code, err := execute(ctx, agent.Start)
	result := CommandResult{
		ID: command.ID,
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

func cmdStop(ctx context.Context, agent Agent, command Command) {
	output, code, err := execute(ctx, agent.Stop)
	result := CommandResult{
		ID: command.ID,
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

func cmdRestart(ctx context.Context, agent Agent, command Command) {
	output, code, err := execute(ctx, agent.Restart)
	result := CommandResult{
		ID: command.ID,
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

func cmdVersion(ctx context.Context, agent Agent, command Command) {
	output, code, err := execute(ctx, agent.Version)
	result := CommandResult{
		ID: command.ID,
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

func cmdStatus(ctx context.Context, agent Agent, command Command) {
	output, code, err := execute(ctx, agent.Status)
	result := CommandResult{
		ID: command.ID,
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
