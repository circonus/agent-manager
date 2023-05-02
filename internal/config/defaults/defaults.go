// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package defaults

import (
	"os"
	"path/filepath"

	"github.com/circonus/go-agent-template/internal/release"
	"github.com/rs/zerolog/log"
)

const (
	// Example arg.
	ExampleArg = "default value"

	// General defaults.

	Debug     = false
	LogLevel  = "info"
	LogPretty = false
)

var (
	// BasePath is the "base" directory
	//
	// expected installation structure:
	// base        (e.g. /opt/circonus/example-agent)
	//   /etc      (e.g. /opt/circonus/example-agent/etc)
	//   /sbin     (e.g. /opt/circonus/example-agent/sbin)
	BasePath = ""

	// EtcPath returns the default etc directory within base directory.
	EtcPath = ""

	// ConfigFile defines the default configuration file name.
	ConfigFile = ""
)

func init() {
	var exePath string
	var resolvedExePath string
	var err error

	exePath, err = os.Executable()
	if err == nil {
		resolvedExePath, err = filepath.EvalSymlinks(exePath)
		if err == nil {
			BasePath = filepath.Clean(filepath.Join(filepath.Dir(resolvedExePath), ".."))
		}
	}

	if err != nil {
		log.Fatal().Err(err).Msg("unable to determine path to binary")
	}

	EtcPath = filepath.Join(BasePath, "etc")
	ConfigFile = filepath.Join(EtcPath, release.NAME+".yaml")
}
