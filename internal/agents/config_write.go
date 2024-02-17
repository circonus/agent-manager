// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import (
	"os"
	"syscall"
)

func writeConfig(path string, data []byte) error {
	perms := os.FileMode(0o640)

	f, err := os.Stat(path)
	if err == nil {
		perms = f.Mode().Perm()
	}

	if err := os.WriteFile(path, data, perms); err != nil {
		return err
	}

	fileSys := f.Sys()
	if s, ok := fileSys.(*syscall.Stat_t); ok {
		if err := os.Chown(path, int(s.Uid), int(s.Gid)); err != nil {
			return err
		}
	}

	return nil
}
