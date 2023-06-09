// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

// Config defines the running configuration options.
type Config struct {
	API             API      `json:"api" toml:"api" yaml:"api"`
	PollingInterval string   `json:"poll_interval" toml:"poll_interval" yaml:"poll_interval"`
	Log             Log      `json:"log" toml:"log" yaml:"log"`
	AWSEC2Tags      []string `json:"aws_ec2_tags"`
	Debug           bool     `json:"debug" toml:"debug" yaml:"debug"`
}

// API defines the various API options.
type API struct {
	URL string `json:"url" toml:"url" yaml:"url"`
}

// Log defines the logging configuration options.
type Log struct {
	Level  string `json:"level" yaml:"level" toml:"level"`
	Pretty bool   `json:"pretty" yaml:"pretty" toml:"pretty"`
}

func Validate() error {
	return nil
}
