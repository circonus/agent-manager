// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	"bytes"
	"os"
)

func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		// docker, while it still works
		return true
	}

	if _, err := os.Stat("/run/.containerenv"); err == nil {
		// podman, while it still works
		return true
	} else if os.Getenv("container") == "podman" {
		return true
	}

	docker := []byte("/docker")
	lxc := []byte("/lxc")

	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil {
		if bytes.Contains(data, docker) || bytes.Contains(data, lxc) {
			return true
		}
	}

	data, err = os.ReadFile("/proc/self/mountinfo")
	if err == nil {
		if bytes.Contains(data, docker) || bytes.Contains(data, lxc) {
			return true
		}
	}

	return false
}
