package collectors

import (
	"context"
	"encoding/base64"
	"os/exec"
	"runtime"
	"strings"

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

func runCommands(ctx context.Context, a Action) error {
	collectors, err := LoadCollectors()
	if err != nil {
		return err
	}
	for _, command := range a.Commands {
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
				args := strings.Split(c.Start, " ")
				cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
				output, err := cmd.CombinedOutput()
				r := Result{
					ActionID: a.ID,
					CommandResult: CommandResult{
						ID: command.ID,
						CommandData: CommandData{
							ExitCode: cmd.ProcessState.ExitCode(),
						},
					},
				}
				if err != nil {
					r.CommandResult.CommandData.Error = err.Error()
				}
				if len(output) > 0 {
					r.CommandResult.CommandData.Output = base64.StdEncoding.EncodeToString(output)
				}
				if err = sendActionResult(ctx, r); err != nil {
					log.Error().Err(err).Msg("command result")
				}
			}
		case STOP:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				args := strings.Split(c.Stop, " ")
				cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
				output, err := cmd.CombinedOutput()
				r := Result{
					ActionID: a.ID,
					CommandResult: CommandResult{
						ID: command.ID,
						CommandData: CommandData{
							ExitCode: cmd.ProcessState.ExitCode(),
						},
					},
				}
				if err != nil {
					r.CommandResult.CommandData.Error = err.Error()
				}
				if len(output) > 0 {
					r.CommandResult.CommandData.Output = base64.StdEncoding.EncodeToString(output)
				}
				if err = sendActionResult(ctx, r); err != nil {
					log.Error().Err(err).Msg("command result")
				}
			}
		case RESTART:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				args := strings.Split(c.Restart, " ")
				cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
				output, err := cmd.CombinedOutput()
				r := Result{
					ActionID: a.ID,
					CommandResult: CommandResult{
						ID: command.ID,
						CommandData: CommandData{
							ExitCode: cmd.ProcessState.ExitCode(),
						},
					},
				}
				if err != nil {
					r.CommandResult.CommandData.Error = err.Error()
				}
				if len(output) > 0 {
					r.CommandResult.CommandData.Output = base64.StdEncoding.EncodeToString(output)
				}
				if err = sendActionResult(ctx, r); err != nil {
					log.Error().Err(err).Msg("command result")
				}
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
				args := strings.Split(c.Status, " ")
				cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
				output, err := cmd.CombinedOutput()
				r := Result{
					ActionID: a.ID,
					CommandResult: CommandResult{
						ID: command.ID,
						CommandData: CommandData{
							ExitCode: cmd.ProcessState.ExitCode(),
						},
					},
				}
				if err != nil {
					r.CommandResult.CommandData.Error = err.Error()
				}
				if len(output) > 0 {
					r.CommandResult.CommandData.Output = base64.StdEncoding.EncodeToString(output)
				}
				if err = sendActionResult(ctx, r); err != nil {
					log.Error().Err(err).Msg("command result")
				}
			}
		case VERSION:
			c, ok := collectors[runtime.GOOS][command.Collector]
			if ok {
				args := strings.Split(c.Version, " ")
				cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
				output, err := cmd.CombinedOutput()
				r := Result{
					ActionID: a.ID,
					CommandResult: CommandResult{
						ID: command.ID,
						CommandData: CommandData{
							ExitCode: cmd.ProcessState.ExitCode(),
						},
					},
				}
				if err != nil {
					r.CommandResult.CommandData.Error = err.Error()
				}
				if len(output) > 0 {
					r.CommandResult.CommandData.Output = base64.StdEncoding.EncodeToString(output)
				}
				if err = sendActionResult(ctx, r); err != nil {
					log.Error().Err(err).Msg("command result")
				}
			}
		}
	}
	return nil
}
