// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package collectors

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
	collectors, err := LoadCollectors()
	if err != nil {
		return err
	}

	for _, command := range action.Commands {
		switch command.Command {
		case INVENTORY:
			if err := FetchCollectors(ctx); err != nil {
				log.Error().Err(err).Msg("refreshing collectors")
			} else if err := CheckForCollectors(ctx); err != nil {
				log.Error().Err(err).Msg("checking for collectors")
			}
		case START:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				cmdStart(ctx, c, command)
			}
		case STOP:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				cmdStop(ctx, c, command)
			}
		case RESTART:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				cmdRestart(ctx, c, command)
			}
		case RELOAD:
			// this needs to be handled differently as reload may be:
			// a command or some type of endpoint
			//
			// c, ok := collectors[runtime.GOOS][command.Collector]
			// if ok {
			// args := strings.Split(c.Reload, " ")
			// cmd := exec.Command(args[0], args[1:]...)
			// }
		case STATUS:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				cmdStatus(ctx, c, command)
			}
		case VERSION:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				cmdVersion(ctx, c, command)
			}
		}
	}

	return nil
}

func cmdStart(ctx context.Context, collector Collector, command Command) {
	output, code, err := execute(ctx, collector.Start)
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

func cmdStop(ctx context.Context, collector Collector, command Command) {
	output, code, err := execute(ctx, collector.Stop)
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

func cmdRestart(ctx context.Context, collector Collector, command Command) {
	output, code, err := execute(ctx, collector.Restart)
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

func cmdVersion(ctx context.Context, collector Collector, command Command) {
	output, code, err := execute(ctx, collector.Version)
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

func cmdStatus(ctx context.Context, collector Collector, command Command) {
	output, code, err := execute(ctx, collector.Status)
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
