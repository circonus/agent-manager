// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//nolint:gochecknoglobals
package release

const (
	// NAME is the name of this application.
	NAME = "circonus-am"
	// ENVPREFIX is the environment variable prefix.
	ENVPREFIX = "CAM"
)

// vars are manipulated at link time (see .goreleaser.yml).
var (
	// COMMIT of release in git repo.
	COMMIT = "none"
	// DATE of release.
	DATE = "unknown"
	// TAG of release.
	TAG = "none"
	// VERSION of the release.
	VERSION = "Dev"
)
