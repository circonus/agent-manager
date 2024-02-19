//go:build linux

// Signal handling for Linux
// doesn't have SIGINFO, using SIGTRAP instead

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
	signal.Notify(m.signalCh, os.Interrupt, unix.SIGTERM, unix.SIGHUP, unix.SIGTRAP)
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
			case unix.SIGTRAP:
				stacklen := runtime.Stack(buf, true)
				fmt.Printf("=== received SIGTRAP ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
			default:
				m.logger.Warn().Str("signal", sig.String()).Msg("unsupported")
			}

		case <-m.groupCtx.Done():
			return nil
		}
	}
}
