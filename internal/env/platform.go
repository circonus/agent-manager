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
