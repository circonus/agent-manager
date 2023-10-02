// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package keys

const (

	// Register - token to use for registering this CMA.
	Register     = "register"
	Inventory    = "inventory"
	Decommission = "decommission"

	APIURL            = "api.url"
	APIToken          = "api.token"
	ManagerID         = "manager_id"
	RegistrationToken = "registration_token"
	RefreshToken      = "refresh_token"
	MachineID         = "machine_id"

	// frequency of polling for actions.
	PollingInterval = "poll_interval"

	// AWS EC2 tags to be included in registration meta data.
	AWSEC2Tags = "aws_ec2_tags"

	// Tags are custom comma separated key:value tags to be added to meta data.
	Tags = "tags"

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

	// UseMachineID - use the machine id or generate a uuid.
	UseMachineID = "use_machine_id"

	// InstanceID - provide an override for ephemeral environments (docker).
	InstanceID = "instance_id"
	// Agents - list of agents manager will manager (docker).
	Agents = "agents"

	//
	// Informational
	// NOTE: Not included in the configuration file, these
	//       options trigger display of information and exit
	//

	// ShowVersion - show version information and exit.
	ShowVersion = "version"

	// Internal settings.
	InventoryFile    = "internal.inventory_file"
	JwtTokenFile     = "internal.jwt_token_file" //nolint:gosec
	ManagerIDFile    = "internal.manager_id_file"
	RefreshTokenFile = "internal.refresh_token_file"
	MachineIDFile    = "internal.machine_id_file"
)
