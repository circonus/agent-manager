package registration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/circonus/collector-management-agent/internal/credentials"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Claims struct {
	Subject string `json:"sub"`
}

// Start the registration process.
func Start(ctx context.Context) error {
	log.Info().Msg("starting registration")

	token := viper.GetString(keys.Register)

	hn, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("getting hostname")
	}
	if hn == "" {
		log.Fatal().Str("hostname", hn).Msg("empty hostname")
	}

	claims := Claims{
		Subject: hn,
	}

	jwt, err := getJWT(ctx, token, claims)
	if err != nil {
		log.Fatal().Err(err).Msg("getting token")
	}

	if err := credentials.Save(jwt); err != nil {
		log.Fatal().Err(err).Msg("saving token")
	}

	return nil
}

func getJWT(ctx context.Context, token string, claims Claims) ([]byte, error) {
	if token == "" {
		return nil, fmt.Errorf("invalid token (empty)") //nolint:goerr113
	}
	if claims.Subject == "" {
		return nil, fmt.Errorf("invalid claims (empty subject)") //nolint:goerr113
	}

	c, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("marshal claims: %w", err)
	}

	reqURL := viper.GetString(keys.APIURL)
	reqURL += "/registration"

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(c))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("X-Circonus-Reg-Token", token)

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

	return body, nil
}
