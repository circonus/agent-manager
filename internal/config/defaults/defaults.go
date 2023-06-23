// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//nolint:gochecknoglobals
package defaults

import (
	"os"
	"path/filepath"

	"github.com/circonus/collector-management-agent/internal/release"
	"github.com/rs/zerolog/log"
)

const (
	APIURL = "https://something.circonus.com"

	PollingInterval = "60s"

	// General defaults.

	Debug     = false
	LogLevel  = "info"
	LogPretty = false
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
	IDPath = ""
	// IDFile is the file where the credentials are stored.
	JwtTokenFile = ""
	RegTokenFile = ""
	AgentIDFile  = ""

	AWSEC2Tags = []string{}
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
	RegTokenFile = filepath.Join(IDPath, "rt")
	AgentIDFile = filepath.Join(IDPath, "ai")

	if err := os.MkdirAll(IDPath, 0700); err != nil {
		log.Fatal().Err(err).Msg("creating ID path")
	}
}
