package credentials

import (
	"fmt"
	"os"

	"github.com/circonus/collector-management-agent/internal/config/defaults"
	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/spf13/viper"
)

func LoadJWT() error {
	if defaults.JwtTokenFile == "" {
		return fmt.Errorf("invalid id file (empty)") //nolint:goerr113
	}

	token, err := os.ReadFile(defaults.JwtTokenFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.APIToken, string(token))

	return nil
}

func SaveJWT(creds []byte) error {
	if defaults.JwtTokenFile == "" {
		return fmt.Errorf("invalid id file (empty)") //nolint:goerr113
	}
	if len(creds) == 0 {
		return fmt.Errorf("invalid credential token (empty)") //nolint:goerr113
	}

	return os.WriteFile(defaults.JwtTokenFile, creds, 0600) //nolint:wrapcheck
}

func LoadAgentID() error {
	if defaults.AgentIDFile == "" {
		return fmt.Errorf("invalid agent id file (empty)") //nolint:goerr113
	}

	token, err := os.ReadFile(defaults.AgentIDFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.AgentID, string(token))

	return nil
}

func SaveAgentID(creds []byte) error {
	if defaults.AgentIDFile == "" {
		return fmt.Errorf("invalid agent id file (empty)") //nolint:goerr113
	}
	if len(creds) == 0 {
		return fmt.Errorf("invalid agent id (empty)") //nolint:goerr113
	}

	return os.WriteFile(defaults.AgentIDFile, creds, 0600) //nolint:wrapcheck
}

func LoadRegistrationToken() error {
	if defaults.RegTokenFile == "" {
		return fmt.Errorf("invalid registration token file (empty)") //nolint:goerr113
	}

	token, err := os.ReadFile(defaults.RegTokenFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.RegistrationToken, string(token))

	return nil
}

func SaveRegistrationToken(creds []byte) error {
	if defaults.RegTokenFile == "" {
		return fmt.Errorf("invalid registration token file (empty)") //nolint:goerr113
	}
	if len(creds) == 0 {
		return fmt.Errorf("invalid registration token (empty)") //nolint:goerr113
	}

	return os.WriteFile(defaults.RegTokenFile, creds, 0600) //nolint:wrapcheck
}
