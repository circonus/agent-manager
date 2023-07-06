package collectors

import (
	"os"
	"syscall"
)

func writeConfig(path string, data []byte) error {
	perms := os.FileMode(0640)

	f, err := os.Stat(path)
	if err == nil {
		perms = f.Mode().Perm()
	}

	if err := os.WriteFile(path, data, perms); err != nil {
		return err //nolint:wrapcheck
	}

	fileSys := f.Sys()
	if _, ok := fileSys.(*syscall.Stat_t); ok {
		gid := int(fileSys.(*syscall.Stat_t).Gid) //nolint:forcetypeassert
		uid := int(fileSys.(*syscall.Stat_t).Uid) //nolint:forcetypeassert

		if err := os.Chown(path, uid, gid); err != nil {
			return err //nolint:wrapcheck
		}
	}

	return nil
}
