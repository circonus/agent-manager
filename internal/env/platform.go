// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package env

import "runtime"

// getPlatform returns the OS and for darwin appends the architecture.
func GetPlatform() string {
	platform := runtime.GOOS
	if platform == "darwin" {
		platform += "_" + runtime.GOARCH
	}

	return platform
}
