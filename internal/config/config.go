// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

// Config defines the running configuration options.
type Config struct {
	Tags            map[string]string `json:"tags"          toml:"tags"          yaml:"tags"`
	API             API               `json:"api"           toml:"api"           yaml:"api"`
	PollingInterval string            `json:"poll_interval" toml:"poll_interval" yaml:"poll_interval"`
	Log             Log               `json:"log"           toml:"log"           yaml:"log"`
	AWSEC2Tags      []string          `json:"aws_ec2_tags"  toml:"aws_ec2_tags"  yaml:"aws_ec2_tags"`
	Debug           bool              `json:"debug"         toml:"debug"         yaml:"debug"`
}

// API defines the various API options.
type API struct {
	URL string `json:"url" toml:"url" yaml:"url"`
}

// Log defines the logging configuration options.
type Log struct {
	Level  string `json:"level"  toml:"level"  yaml:"level"`
	Pretty bool   `json:"pretty" toml:"pretty" yaml:"pretty"`
}

func Validate() error {
	return nil
}
