// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

// Config defines the running configuration options.
type Config struct {
	ExampleArg string `json:"example_arg" toml:"example_arg" yaml:"example_arg"`
	Log        Log    `json:"log" toml:"log" yaml:"log"`
	Debug      bool   `json:"debug" toml:"debug" yaml:"debug"`
}

// Log defines the logging configuration options.
type Log struct {
	Level  string `json:"level" yaml:"level" toml:"level"`
	Pretty bool   `json:"pretty" yaml:"pretty" toml:"pretty"`
}

func Validate() error {
	return nil
}
