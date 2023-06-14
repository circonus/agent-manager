package collectors

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"

	"github.com/circonus/collector-management-agent/internal/config/defaults"
	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// handle requesting list of collectors from api, determining if any are installed locally, and responding with collectors found

// key1 is platform (GOOS) e.g. darwin, windows, linux, freebsd, etc.
// key2 is collector e.g. fluent-bit, telegraf, etc.
type Collectors map[string]map[string]Collector

type Collector struct {
	Binary      string   `json:"binary" yaml:"binary"`
	Start       string   `json:"start" yaml:"start"`
	Stop        string   `json:"stop" yaml:"stop"`
	Restart     string   `json:"restart" yaml:"restart"`
	Reload      string   `json:"reload" yaml:"reload"`
	Status      string   `json:"status" yaml:"status"`
	ConfigFiles []string `json:"config_files" yaml:"config_files"`
}

func FetchCollectors(ctx context.Context) error {
	// /collector_type for list of known collectors
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "collector_type")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("X-Circonus-Token", token)

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response -- status: %s, body: %s", resp.Status, string(body)) //nolint:goerr113
	}

	var collectors Collectors
	if err := json.Unmarshal(body, &collectors); err != nil {
		return fmt.Errorf("unmarshal body: %w", err)
	}

	return SaveCollectors(collectors)
}

func LoadCollectors() (Collectors, error) {
	data, err := os.ReadFile(defaults.InvetoryFile)
	if err != nil {
		return nil, fmt.Errorf("loading collector inventory: %w", err)
	}

	var c Collectors
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing collector inventory: %w", err)
	}

	return c, nil
}

func SaveCollectors(c Collectors) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal collector inventory: %w", err)
	}

	if err := os.WriteFile(defaults.InvetoryFile, data, 0600); err != nil {
		return fmt.Errorf("saving collector inventory: %w", err)
	}

	return nil
}

func CheckForCollectors(ctx context.Context) error {
	cc, err := LoadCollectors()
	if err != nil {
		return err
	}

	gcc, ok := cc[runtime.GOOS]
	if !ok {
		return fmt.Errorf("no collectors found for platform %s", runtime.GOOS) //nolint:goerr113
	}

	found := []string{}

	for name, c := range gcc {
		if _, err := os.Stat(c.Binary); errors.Is(err, os.ErrNotExist) {
			log.Warn().Str("file", c.Binary).Msg("collector binary not found, skipping")
			continue
		}
		for _, config := range c.ConfigFiles {
			if _, err := os.Stat(config); errors.Is(err, os.ErrNotExist) {
				log.Warn().Str("file", config).Msg("required config file not found, skipping")
				continue
			}
		}

		log.Info().Str("collector agent", name).Msg("found")
		found = append(found, name)
	}

	if len(found) > 0 {
		// contact api and report what collectors were found
		return registerCollectors(ctx, found)
	}

	return nil
}

func registerCollectors(ctx context.Context, c []string) error {
	token := viper.GetString(keys.APIToken)
	if token == "" {
		return fmt.Errorf("invalid api token (empty)") //nolint:goerr113
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "collector")
	if err != nil {
		return fmt.Errorf("req url: %w", err)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal claims: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("X-Circonus-Token", token)

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response -- status: %s, body: %s", resp.Status, string(body)) //nolint:goerr113
	}

	return nil
}
