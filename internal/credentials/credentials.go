package credentials

import (
	"fmt"
	"os"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/spf13/viper"
)

func LoadJWT() error {
	file := viper.GetString(keys.JwtTokenFile)
	if file == "" {
		return fmt.Errorf("invalid id file (empty)") //nolint:goerr113
	}

	token, err := read(file)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.APIToken, string(token))

	return nil
}

func SaveJWT(creds []byte) error {
	file := viper.GetString(keys.JwtTokenFile)
	if file == "" {
		return fmt.Errorf("invalid id file (empty)") //nolint:goerr113
	}
	if len(creds) == 0 {
		return fmt.Errorf("invalid credential token (empty)") //nolint:goerr113
	}

	return write(file, creds) //nolint:wrapcheck
}

func LoadAgentID() error {
	file := viper.GetString(keys.AgentIDFile)
	if file == "" {
		return fmt.Errorf("invalid agent id file (empty)") //nolint:goerr113
	}

	token, err := read(file)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.AgentID, string(token))

	return nil
}

func SaveAgentID(creds []byte) error {
	file := viper.GetString(keys.AgentIDFile)
	if file == "" {
		return fmt.Errorf("invalid agent id file (empty)") //nolint:goerr113
	}
	if len(creds) == 0 {
		return fmt.Errorf("invalid agent id (empty)") //nolint:goerr113
	}

	return write(file, creds) //nolint:wrapcheck
}

func LoadRegistrationToken() error {
	file := viper.GetString(keys.RegTokenFile)
	if file == "" {
		return fmt.Errorf("invalid registration token file (empty)") //nolint:goerr113
	}

	token, err := read(file)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.RegistrationToken, string(token))

	return nil
}

func SaveRegistrationToken(creds []byte) error {
	file := viper.GetString(keys.RegTokenFile)
	if file == "" {
		return fmt.Errorf("invalid registration token file (empty)") //nolint:goerr113
	}
	if len(creds) == 0 {
		return fmt.Errorf("invalid registration token (empty)") //nolint:goerr113
	}

	return write(file, creds) //nolint:wrapcheck
}

func read(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return data, nil
}

func write(file string, data []byte) error {
	return os.WriteFile(file, data, 0600) //nolint:wrapcheck
}
