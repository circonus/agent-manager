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
	"time"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/circonus/collector-management-agent/internal/credentials"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// re-register every hour to get new JWT

func ReRegister(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			reg, err := getNewJWT(ctx)
			if err != nil {
				log.Error().Err(err).Msg("getting token")
			}

			if err := credentials.SaveJWT([]byte(reg.AuthToken)); err != nil {
				log.Fatal().Err(err).Msg("saving token")
			}
		}
	}
}

func getNewJWT(ctx context.Context) (*Response, error) {
	token := viper.GetString(keys.RegistrationToken)
	if token == "" {
		return nil, fmt.Errorf("invalid token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent", "register")
	if err != nil {
		return nil, fmt.Errorf("req url: %w", err)
	}

	data := []byte(`{"agent_id":"` + viper.GetString(keys.AgentID) + `"}`)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("X-Circonus-Register-Token", token)

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
		return nil, fmt.Errorf("non-200 response -- status: %s, body: %s", resp.Status, string(body)) //nolint:goerr113
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parsing response body: %w", err)
	}

	return &response, nil
}
