// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/credentials"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// refresh registration every hour to get new JWT

func RefreshRegistration(ctx context.Context) error {
	log.Info().Msg("refreshing token")

	reg, err := getNewJWT(ctx)
	if err != nil {
		return fmt.Errorf("refreshing token: %w", err)
	}

	if err := credentials.SaveRefreshToken([]byte(reg.RefreshToken)); err != nil {
		return fmt.Errorf("saving refresh token: %w", err)
	}

	if err := credentials.SaveJWT([]byte(reg.AuthToken)); err != nil {
		return fmt.Errorf("saving access token: %w", err)
	}

	if err := credentials.LoadJWT(); err != nil {
		return fmt.Errorf("loading access token: %w", err)
	}

	return nil
}

func getNewJWT(ctx context.Context) (*Response, error) {
	if err := credentials.LoadRefreshToken(); err != nil {
		return nil, fmt.Errorf("loading refresh token: %w", err)
	}

	token := viper.GetString(keys.RefreshToken)
	if token == "" {
		return nil, fmt.Errorf("invalid refresh token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "manager", "register")
	if err != nil {
		return nil, fmt.Errorf("req url: %w", err)
	}

	data := []byte(`{"manager_id":"` + viper.GetString(keys.ManagerID) + `"}`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling registration endpoint: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response -- status: %d %s, body: %s", resp.StatusCode, resp.Status, string(body)) //nolint:goerr113
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parsing response body: %w", err)
	}

	return &response, nil
}
