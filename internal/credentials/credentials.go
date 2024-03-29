package credentials

import (
	"errors"
	"fmt"
	"os"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/spf13/viper"
)

func LoadJWT() error {
	file := viper.GetString(keys.JwtTokenFile)
	if file == "" {
		return fmt.Errorf("invalid id file (empty)")
	}

	token, err := read(file)
	if err != nil {
		return err
	}

	viper.Set(keys.APIToken, string(token))

	return nil
}

func SaveJWT(creds []byte) error {
	file := viper.GetString(keys.JwtTokenFile)
	if file == "" {
		return fmt.Errorf("invalid id file (empty)")
	}

	if len(creds) == 0 {
		return fmt.Errorf("invalid credential token (empty)")
	}

	return write(file, creds)
}

func LoadManagerID() error {
	file := viper.GetString(keys.ManagerIDFile)
	if file == "" {
		return fmt.Errorf("invalid manager id file (empty)")
	}

	token, err := read(file)
	if err != nil {
		return err
	}

	viper.Set(keys.ManagerID, string(token))

	return nil
}

func SaveManagerID(creds []byte) error {
	file := viper.GetString(keys.ManagerIDFile)
	if file == "" {
		return fmt.Errorf("invalid manager id file (empty)")
	}

	if len(creds) == 0 {
		return fmt.Errorf("invalid manager id (empty)")
	}

	return write(file, creds)
}

func LoadRefreshToken() error {
	file := viper.GetString(keys.RefreshTokenFile)
	if file == "" {
		return fmt.Errorf("invalid refresh token file (empty)")
	}

	token, err := read(file)
	if err != nil {
		return err
	}

	viper.Set(keys.RefreshToken, string(token))

	return nil
}

func SaveRefreshToken(creds []byte) error {
	file := viper.GetString(keys.RefreshTokenFile)
	if file == "" {
		return fmt.Errorf("invalid refresh token file (empty)")
	}

	if len(creds) == 0 {
		return fmt.Errorf("invalid refresh token (empty)")
	}

	return write(file, creds)
}

func LoadMachineID() error {
	file := viper.GetString(keys.MachineIDFile)
	if file == "" {
		return fmt.Errorf("invalid machine id file (empty)")
	}

	token, err := read(file)
	if err != nil {
		return err
	}

	viper.Set(keys.MachineID, string(token))

	return nil
}

func SaveMachineID(creds []byte) error {
	file := viper.GetString(keys.MachineIDFile)
	if file == "" {
		return fmt.Errorf("invalid machine id file (empty)")
	}

	if len(creds) == 0 {
		return fmt.Errorf("invalid machine id (empty)")
	}

	return write(file, creds)
}

func read(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func write(file string, data []byte) error {
	return os.WriteFile(file, data, 0o600)
}

func DoesFileExist(file string) bool {
	if fs, err := os.Stat(file); err == nil {
		if fs.Size() > 0 {
			return true
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false
}
