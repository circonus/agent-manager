// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package decomission

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/credentials"
	"github.com/spf13/viper"
)

func Start(ctx context.Context) error {
	if err := credentials.LoadManagerID(); err != nil {
		return fmt.Errorf("loading manager id: %w", err)
	}

	if err := credentials.LoadJWT(); err != nil {
		return fmt.Errorf("loading API credentials: %w", err)
	}

	if err := deleteManager(ctx); err != nil {
		return fmt.Errorf("deleting agent manager record: %w", err)
	}

	if err := os.Remove(viper.GetString(keys.InventoryFile)); err != nil {
		return fmt.Errorf("removing %s: %w", viper.GetString(keys.InventoryFile), err)
	}

	if err := os.Remove(viper.GetString(keys.JwtTokenFile)); err != nil {
		return fmt.Errorf("removing %s: %w", viper.GetString(keys.JwtTokenFile), err)
	}

	if err := os.Remove(viper.GetString(keys.ManagerIDFile)); err != nil {
		return fmt.Errorf("removing %s: %w", viper.GetString(keys.ManagerIDFile), err)
	}

	if err := os.Remove(viper.GetString(keys.RefreshTokenFile)); err != nil {
		return fmt.Errorf("removing %s: %w", viper.GetString(keys.RefreshTokenFile), err)
	}

	return nil
}

func deleteManager(ctx context.Context) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "manager", viper.GetString(keys.ManagerID))
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("calling actions endpoint: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("non-204 response -- status: %s, body: %s", resp.Status, string(body)) //nolint:goerr113
	}

	return nil
}
