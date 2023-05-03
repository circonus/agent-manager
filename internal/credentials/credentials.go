package credentials

import (
	"fmt"
	"os"

	"github.com/circonus/collector-management-agent/internal/config/defaults"
	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/spf13/viper"
)

func Load() error {
	if defaults.IDFile == "" {
		return fmt.Errorf("invalid id file (empty)") //nolint:goerr113
	}

	token, err := os.ReadFile(defaults.IDFile)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.APIToken, string(token))

	return nil
}

func Save(creds []byte) error {
	if defaults.IDFile == "" {
		return fmt.Errorf("invalid id file (empty)") //nolint:goerr113
	}
	if len(creds) == 0 {
		return fmt.Errorf("invalid credential token (empty)") //nolint:goerr113
	}

	return os.WriteFile(defaults.IDFile, creds, 0600) //nolint:wrapcheck
}
