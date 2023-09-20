// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

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

func LoadManagerID() error {
	file := viper.GetString(keys.ManagerIDFile)
	if file == "" {
		return fmt.Errorf("invalid manager id file (empty)") //nolint:goerr113
	}

	token, err := read(file)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.ManagerID, string(token))

	return nil
}

func SaveManagerID(creds []byte) error {
	file := viper.GetString(keys.ManagerIDFile)
	if file == "" {
		return fmt.Errorf("invalid manager id file (empty)") //nolint:goerr113
	}

	if len(creds) == 0 {
		return fmt.Errorf("invalid manager id (empty)") //nolint:goerr113
	}

	return write(file, creds) //nolint:wrapcheck
}

func LoadRefreshToken() error {
	file := viper.GetString(keys.RefreshTokenFile)
	if file == "" {
		return fmt.Errorf("invalid refresh token file (empty)") //nolint:goerr113
	}

	token, err := read(file)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.RefreshToken, string(token))

	return nil
}

func SaveRefreshToken(creds []byte) error {
	file := viper.GetString(keys.RefreshTokenFile)
	if file == "" {
		return fmt.Errorf("invalid refresh token file (empty)") //nolint:goerr113
	}

	if len(creds) == 0 {
		return fmt.Errorf("invalid refresh token (empty)") //nolint:goerr113
	}

	return write(file, creds) //nolint:wrapcheck
}

func LoadMachineID() error {
	file := viper.GetString(keys.MachineID)
	if file == "" {
		return fmt.Errorf("invalid machine id file (empty)") //nolint:goerr113
	}

	token, err := read(file)
	if err != nil {
		return err //nolint:wrapcheck
	}

	viper.Set(keys.MachineIDFile, string(token))

	return nil
}

func SaveMachineID(creds []byte) error {
	file := viper.GetString(keys.MachineIDFile)
	if file == "" {
		return fmt.Errorf("invalid machine id file (empty)") //nolint:goerr113
	}

	if len(creds) == 0 {
		return fmt.Errorf("invalid machine id (empty)") //nolint:goerr113
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

func DoesFileExist(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false
}
