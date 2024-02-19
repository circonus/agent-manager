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
