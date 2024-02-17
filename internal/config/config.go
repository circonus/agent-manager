// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//nolint:lll
package config

import (
	"github.com/circonus/agent-manager/internal/config/defaults"
)

// Config defines the running configuration options.
type Config struct {
	Tags                   map[string]string `json:"tags"                  toml:"tags"                  yaml:"tags"`
	API                    API               `json:"api"                   toml:"api"                   yaml:"api"`
	ActionPollingInterval  string            `json:"action_poll_interval"  toml:"action_poll_interval"  yaml:"action_poll_interval"`
	TrackerPollingInterval string            `json:"tracker_poll_interval" toml:"tracker_poll_interval" yaml:"tracker_poll_interval"`
	StatusPollingInterval  string            `json:"status_poll_interval"  toml:"status_poll_interval"  yaml:"status_poll_interval"`
	Server                 Server            `json:"server"                toml:"server"                yaml:"server"`
	Log                    Log               `json:"log"                   toml:"log"                   yaml:"log"`
	AWSEC2Tags             []string          `json:"aws_ec2_tags"          toml:"aws_ec2_tags"          yaml:"aws_ec2_tags"`
	Debug                  bool              `json:"debug"                 toml:"debug"                 yaml:"debug"`
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

type Server struct {
	Address           string `json:"address"             toml:"address"             yaml:"address"`
	ReadTimeout       string `json:"read_timeout"        toml:"read_timeout"        yaml:"read_timeout"`
	WriteTimeout      string `json:"write_timeout"       toml:"write_timeout"       yaml:"write_timeout"`
	IdleTimeout       string `json:"idle_timeout"        toml:"idle_timeout"        yaml:"idle_timeout"`
	ReadHeaderTimeout string `json:"read_header_timeout" toml:"read_header_timeout" yaml:"read_header_timeout"`
	HandlerTimeout    string `json:"handler_timeout"     toml:"handler_timeout"     yaml:"handler_timeout"`
	TLSKeyFile        string `json:"tls_key_file"        toml:"tls_key_file"        yaml:"tls_key_file"`
	TLSCertFile       string `json:"tls_cert_file"       toml:"tls_cert_file"       yaml:"tls_cert_file"`
	TLSEnable         bool   `json:"tls_enable"          toml:"tls_enable"          yaml:"tls_enable"`
}

func Validate() error {
	return nil
}

func SetPathsBasedOnConfigFile(cfgPath string) {
	if cfgPath == "" {
		return
	}

	if cfgPath == defaults.EtcPath {
		return
	}

	defaults.SetEtcPaths(cfgPath)
}
