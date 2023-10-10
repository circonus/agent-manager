// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package inventory

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/rs/zerolog/log"
)

func backupConfigs(name string, configs map[string]string) {
	baseDir := filepath.Join(defaults.EtcPath, "configs", name)
	// e.g. /opt/circonus/am/etc/configs/telegraf

	if err := os.MkdirAll(baseDir, 0700); err != nil {
		if errors.Is(err, os.ErrExist) {
			return // don't need to backup again
		}

		if !errors.Is(err, os.ErrExist) {
			log.Error().Err(err).Str("path", baseDir).Msg("unable to make config dir to save backup")

			return
		}
	}

	for _, src := range configs {
		sfi, err := os.Stat(src)
		if err != nil {
			log.Error().Err(err).Str("src", src).Msg("stat source file")

			return
		}

		if !sfi.Mode().IsRegular() {
			log.Error().Str("src", src).Str("mode", sfi.Mode().String()).Msg("source is not a regular file")

			return
		}

		dst := filepath.Join(baseDir, filepath.Base(src))

		in, err := os.Open(src)
		if err != nil {
			log.Error().Err(err).Str("src", src).Msg("opening source file")

			return
		}

		out, err := os.Create(dst)
		if err != nil {
			log.Error().Err(err).Str("dst", dst).Msg("creating destination file")
			in.Close()

			return
		}

		if _, err := io.Copy(out, in); err != nil {
			log.Error().Err(err).Msg("copying file contents")
			in.Close()
			out.Close()

			return
		}

		in.Close()
		out.Close()
		log.Info().Str("src", src).Str("dst", dst).Msg("backed up current config file")
	}
}
