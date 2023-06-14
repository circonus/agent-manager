// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package keys

const (

	// Register - token to use for registering this CMA.
	Register = "register"

	APIURL   = "api.url"
	APIToken = "api.token"

	// frequency of polling for actions.
	PollingInterval = "poll_interval"

	// AWS EC2 tags to be included in registration meta data.
	AWSEC2Tags = "aws_ec2_tags"

	//
	// Logging.
	//

	// LogLevel logging level (panic, fatal, error, warn, info, debug, disabled).
	LogLevel = "log.level"

	// LogPretty output formatted log lines (for running in foreground).
	LogPretty = "log.pretty"

	//
	// Miscellaneous.
	//

	// Debug enables debug messages.
	Debug = "debug"

	//
	// Informational
	// NOTE: Not included in the configuration file, these
	//       options trigger display of information and exit
	//

	// ShowVersion - show version information and exit.
	ShowVersion = "version"
)
