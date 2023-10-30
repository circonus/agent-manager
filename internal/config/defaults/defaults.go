// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//nolint:gochecknoglobals
package defaults

import (
	"os"
	"path/filepath"

	"github.com/circonus/agent-manager/internal/release"
	"github.com/rs/zerolog/log"
)

const (
	APIURL = "https://agents-api.circonus.app/configurations/v1"

	ActionPollingInterval  = "60s"
	TrackerPollingInterval = "15m"
	StatusPollingInterval  = "5m"

	// General defaults.

	Debug     = false
	LogLevel  = "info"
	LogPretty = false

	UseMachineID  = true
	ForceRegister = false

	ServerAddress           = ":43285"
	ServerReadTimeout       = "60s"
	ServerWriteTimeout      = "60s"
	ServerIdleTimeout       = "30s"
	ServerReadHeaderTimeout = "5s"
	ServerHandlerTimeout    = "30s"
	ServerUseTLS            = false
	ServerCertFile          = ""
	ServerKeyFile           = ""
)

var (
	// BasePath is the "base" directory
	//
	// expected installation structure:
	// base        (e.g. /opt/circonus/cma)
	//   /etc      (e.g. /opt/circonus/cma/etc)
	//      /.id   (e.g. /opt/circonus/cma/etc/.id)
	//   /sbin     (e.g. /opt/circonus/cma/sbin)
	BasePath = ""

	// EtcPath returns the default etc directory within base directory.
	EtcPath = ""

	// ConfigFile defines the default configuration file name.
	ConfigFile = ""

	// Collector inventory file.
	InventoryFile = ""

	// IDPath is where ID credentials are stored.
	IDPath           = ""
	JwtTokenFile     = ""
	ManagerIDFile    = ""
	RefreshTokenFile = ""
	MachineIDFile    = ""

	AWSEC2Tags = []string{}
	Tags       = []string{}
	Agents     = []string{}
)

func init() { //nolint:gochecknoinits
	var exePath, resolvedExePath string

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

	// Set default etc paths/files. Will be re-set if a specific file was
	// identified via --config argument.
	SetEtcPaths(filepath.Join(BasePath, "etc"))
}

func SetEtcPaths(etcPath string) {
	EtcPath = etcPath
	ConfigFile = filepath.Join(EtcPath, release.NAME+".yaml")
	InventoryFile = filepath.Join(EtcPath, "inventory.yaml")
	IDPath = filepath.Join(EtcPath, ".id")
	JwtTokenFile = filepath.Join(IDPath, "jt")
	ManagerIDFile = filepath.Join(IDPath, "ai")
	RefreshTokenFile = filepath.Join(IDPath, "rft")
	MachineIDFile = filepath.Join(IDPath, "mid")

	if err := os.MkdirAll(IDPath, 0700); err != nil {
		log.Fatal().Err(err).Msg("creating ID path")
	}
}
