package registration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/credentials"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/spf13/viper"
)

func UpdateVersion(ctx context.Context) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		if err := credentials.LoadJWT(); err != nil {
			return fmt.Errorf("loading api token: %w", err)
		}

		token = viper.GetString(keys.APIToken)
		if token == "" {
			return fmt.Errorf("invalid api token (empty)")
		}
	}

	managerID := viper.GetString(keys.ManagerID)
	if managerID == "" {
		if err := credentials.LoadManagerID(); err != nil {
			return fmt.Errorf("loading manager id: %w", err)
		}

		managerID = viper.GetString(keys.ManagerID)
		if managerID == "" {
			return fmt.Errorf("invalid manager ID (empty)")
		}
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "manager", managerID)
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	data := []byte(`{"version":"v` + release.VERSION + `"}`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("calling registration endpoint: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response -- status: %d %s, body: %s", resp.StatusCode, resp.Status, string(body))
	}

	return nil
}
