package main

import (
	"github.com/circonus/agent-manager/internal/config/defaults"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initAppArgs adds application specific args to the cobra command.
func initAppArgs(cmd *cobra.Command) {
	{
		const (
			key          = keys.Register
			longOpt      = "register"
			description  = "Registration token -- register agent manager, inventory installed agents and exit"
			envVar       = release.ENVPREFIX + "_REGISTER"
			defaultValue = ""
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}
	{
		const (
			key          = keys.ForceRegister
			longOpt      = "force-register"
			envVar       = release.ENVPREFIX + "_FORCE_REGISTER"
			description  = "Force registration attempt, even if manager is already registered"
			defaultValue = defaults.ForceRegister
		)

		cmd.Flags().Bool(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.Decommission
			longOpt      = "decommission"
			description  = "Decommission agent manager and exit"
			defaultValue = false
		)

		cmd.Flags().Bool(longOpt, defaultValue, description)
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.APIURL
			longOpt      = "apiurl"
			envVar       = release.ENVPREFIX + "_API_URL"
			description  = "Circonus API URL"
			defaultValue = defaults.APIURL
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ActionPollingInterval
			longOpt      = "action-poll-interval"
			envVar       = release.ENVPREFIX + "_ACTION_POLL_INTERVAL"
			description  = "Polling interval for actions"
			defaultValue = defaults.ActionPollingInterval
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.TrackerPollingInterval
			longOpt      = "tracker-poll-interval"
			envVar       = release.ENVPREFIX + "_TRACKER_POLL_INTERVAL"
			description  = "Polling interval for tracking and verifying checksums"
			defaultValue = defaults.TrackerPollingInterval
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.StatusPollingInterval
			longOpt      = "status-poll-interval"
			envVar       = release.ENVPREFIX + "_STATUS_POLL_INTERVAL"
			description  = "Polling interval for gathering agent status"
			defaultValue = defaults.StatusPollingInterval
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = keys.AWSEC2Tags
			longOpt     = "aws-ec2-tags"
			envVar      = release.ENVPREFIX + "_AWS_EC2_TAGS"
			description = "AWS EC2 tags for registration meta data"
		)

		defaultValue := defaults.AWSEC2Tags

		cmd.Flags().StringSlice(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = keys.Tags
			longOpt     = "tags"
			envVar      = release.ENVPREFIX + "_TAGS"
			description = "Custom key:value tags for registration meta data"
			// env separate with space CAM_TAGS="foo:bar baz:qux"
			// cli separate with comma --tags="foo:bar,baz:qux"
		)

		defaultValue := defaults.Tags

		cmd.Flags().StringSlice(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = keys.UseMachineID
			longOpt     = "use-machine-id"
			envVar      = release.ENVPREFIX + "_USE_MACHINE_ID"
			description = "Use machine_id or generate uuid"
		)

		defaultValue := defaults.UseMachineID

		cmd.Flags().Bool(longOpt, defaultValue, envDescription(description, envVar))
		flag := cmd.Flags().Lookup(longOpt)
		flag.Hidden = true
		bindFlagError(longOpt, viper.BindPFlag(key, flag))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.InstanceID
			longOpt      = "instance-id"
			envVar       = release.ENVPREFIX + "_INSTANCE_ID"
			description  = "Instance ID (Docker specific)"
			defaultValue = ""
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = keys.Agents
			longOpt     = "agents"
			envVar      = release.ENVPREFIX + "_AGENTS"
			description = "List of agents (Docker specific)"
		)

		defaultValue := defaults.Agents

		cmd.Flags().StringSlice(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerAddress
			longOpt      = "server-address"
			envVar       = release.ENVPREFIX + "_SERVER_ADDRESS"
			description  = "Server Address for /health and /config"
			defaultValue = defaults.ServerAddress
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerReadTimeout
			longOpt      = "server-read-timeout"
			envVar       = release.ENVPREFIX + "_SERVER_READ_TIMEOUT"
			description  = "Server read timeout"
			defaultValue = defaults.ServerReadTimeout
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerWriteTimeout
			longOpt      = "server-write-timeout"
			envVar       = release.ENVPREFIX + "_SERVER_WRITE_TIMEOUT"
			description  = "Server write timeout"
			defaultValue = defaults.ServerWriteTimeout
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerIdleTimeout
			longOpt      = "server-idle-timeout"
			envVar       = release.ENVPREFIX + "_SERVER_IDLE_TIMEOUT"
			description  = "Server idle timeout"
			defaultValue = defaults.ServerIdleTimeout
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerReadHeaderTimeout
			longOpt      = "server-read-header-timeout"
			envVar       = release.ENVPREFIX + "_SERVER_READ_HEADER_TIMEOUT"
			description  = "Server read header timeout"
			defaultValue = defaults.ServerReadHeaderTimeout
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerHandlerTimeout
			longOpt      = "server-handler-timeout"
			envVar       = release.ENVPREFIX + "_SERVER_HANDLER_TIMEOUT"
			description  = "Server handler timeout"
			defaultValue = defaults.ServerHandlerTimeout
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = keys.ServerTLSEnable
			longOpt     = "server-tls-enable"
			envVar      = release.ENVPREFIX + "_SERVER_TLS_ENABLE"
			description = "Server Enable TLS"
		)

		defaultValue := defaults.ServerUseTLS

		cmd.Flags().Bool(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerTLSKeyFile
			longOpt      = "server-tls-key-file"
			envVar       = release.ENVPREFIX + "_SERVER_TLS_KEY_FILE"
			description  = "Server TLS key file"
			defaultValue = defaults.ServerKeyFile
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = keys.ServerTLSCertFile
			longOpt      = "server-tls-cert-file"
			envVar       = release.ENVPREFIX + "_SERVER_TLS_CERT_FILE"
			description  = "Server TLS cert file"
			defaultValue = defaults.ServerKeyFile
		)

		cmd.Flags().String(longOpt, defaultValue, envDescription(description, envVar))
		bindFlagError(longOpt, viper.BindPFlag(key, cmd.Flags().Lookup(longOpt)))
		bindEnvError(envVar, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}
}
