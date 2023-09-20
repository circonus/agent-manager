// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

import "runtime"

// getPlatform returns the OS and for darwin appends the architecture.
func getPlatform() string {
	platform := runtime.GOOS
	if platform == "darwin" {
		platform += "_" + runtime.GOARCH
	}

	return platform
}
