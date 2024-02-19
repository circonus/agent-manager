package server

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type healthHandler struct{}

func (healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:       5 * time.Second,
				KeepAlive:     3 * time.Second,
				FallbackDelay: -1 * time.Millisecond,
			}).DialContext,
			DisableKeepAlives:   true,
			DisableCompression:  false,
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 0,
		},
		Timeout: 10 * time.Second,
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "health")
	if err != nil {
		log.Error().Err(err).Str("api_url", viper.GetString(keys.APIURL)).Msg("creating API health URL")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		log.Error().Err(err).Str("api_url", viper.GetString(keys.APIURL)).Msg("creating API health request")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "close")
	req.Close = true

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("api_url", viper.GetString(keys.APIURL)).Msg("requesting API health URL")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Str("api_url", viper.GetString(keys.APIURL)).Msg("reading response from API health URL")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Err(err).
			Str("api_url", viper.GetString(keys.APIURL)).
			Int("status_code", resp.StatusCode).
			Str("status", resp.Status).
			Str("body", string(body)).
			Msg("requesting API health URL")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	_, _ = w.Write([]byte(`{"status":"ok","dur":"` + time.Since(start).String() + `"}`))
}
