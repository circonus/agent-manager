//go:build windows

// Signal handling for Windows
// doesn't have SIGINFO, attempt to use SIGTRAP instead...

package manager

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/alecthomas/units"
)

func (m *Manager) signalNotifySetup() {
	signal.Notify(m.signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTRAP)
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
			case os.Interrupt, syscall.SIGTERM:
				m.Stop()
			case syscall.SIGHUP:
				// Noop
			case syscall.SIGTRAP:
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
