// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build freebsd || openbsd || solaris || darwin

// Signal handling for FreeBSD, OpenBSD, Darwin, and Solaris
// systems that have SIGINFO

package manager

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/alecthomas/units"
	"golang.org/x/sys/unix"
)

func (m *Manager) signalNotifySetup() {
	signal.Notify(m.signalCh, os.Interrupt, unix.SIGTERM, unix.SIGHUP, unix.SIGINFO)
}

// handleSignals runs the signal handler thread.
func (m *Manager) handleSignals() error {
	const stacktraceBufSize = 1 * units.MiB

	// pre-allocate a buffer
	buf := make([]byte, stacktraceBufSize)

	for {
		select {
		case sig := <-m.signalCh:
			m.logger.Info().Str("signal", sig.String()).Msg("received signal")

			switch sig {
			case os.Interrupt, unix.SIGTERM:
				m.Stop()
			case unix.SIGHUP:
				// Noop
			case unix.SIGINFO:
				stacklen := runtime.Stack(buf, true)
				fmt.Printf("=== received SIGINFO ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
			default:
				m.logger.Warn().Str("signal", sig.String()).Msg("unsupported")
			}

		case <-m.groupCtx.Done():
			return nil
		}
	}
}
