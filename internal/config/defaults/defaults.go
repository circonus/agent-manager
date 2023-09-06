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
	APIURL = "https://web-api.svcs-np.circonus.net/configurations/v1"

	PollingInterval = "60s"

	// General defaults.

	Debug     = false
	LogLevel  = "info"
	LogPretty = false

	UseMachineID = true
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

	EtcPath = filepath.Join(BasePath, "etc")
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
