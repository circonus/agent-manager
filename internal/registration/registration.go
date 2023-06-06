package registration

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/circonus/collector-management-agent/internal/credentials"
	"github.com/denisbrodbeck/machineid"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/viper"
)

type Request struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Hostname             string `json:"hostname"`
	MachineID            string `json:"machine_id"`
	OS                   string `json:"os"`
	Platform             string `json:"platform"`
	PlatformVersion      string `json:"platform_version"`
	PlatformFamily       string `json:"platform_family"`
	KernelVersion        string `json:"kernel_version"`
	KernelArch           string `json:"kernel_arch"`
	VirtualizationSystem string `json:"virtualization_system"`
	VirtualizationRole   string `json:"virtualization_role"`
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

	mid, err := getMachineID()
	if err != nil {
		log.Fatal().Err(err).Str("mid", mid).Msg("invalid machine id")
	}

	meta, err := getHostInfo()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to retrieve host info")
	}

	meta.Hostname = hn
	meta.MachineID = mid

	request := Request{
		Meta: meta,
	}

	jwt, err := getJWT(ctx, token, request)
	if err != nil {
		log.Fatal().Err(err).Msg("getting token")
	}

	if err := credentials.Save(jwt); err != nil {
		log.Fatal().Err(err).Msg("saving token")
	}

	return nil
}

func getJWT(ctx context.Context, token string, request Request) ([]byte, error) {
	if token == "" {
		return nil, fmt.Errorf("invalid token (empty)") //nolint:goerr113
	}
	if request.Meta.Hostname == "" {
		return nil, fmt.Errorf("invalid claims (empty hostname)") //nolint:goerr113
	}
	if request.Meta.MachineID == "" {
		return nil, fmt.Errorf("invalid claims (empty machine id)") //nolint:goerr113
	}

	c, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal claims: %w", err)
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "registration")
	if err != nil {
		return nil, fmt.Errorf("req url: %w", err)
	}

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

func getMachineID() (string, error) {
	id, err := machineid.ID()
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	mac := hmac.New(sha256.New, []byte(id))

	return fmt.Sprintf("%x", mac.Sum(nil)), nil
}

func getHostInfo() (Meta, error) {

	meta := Meta{}

	hi, err := host.Info()
	if err != nil {
		return meta, err //nolint:wrapcheck
	}

	// hi.OS is runtime.GOOS
	if hi.OS != "" {
		meta.OS = hi.OS
	}

	if hi.Platform != "" {
		meta.Platform = hi.Platform
	}
	if hi.PlatformFamily != "" {
		meta.PlatformFamily = hi.PlatformFamily
	}
	if hi.PlatformVersion != "" {
		meta.PlatformVersion = hi.PlatformVersion
	}
	if hi.KernelVersion != "" {
		meta.KernelVersion = hi.KernelVersion
	}
	if hi.KernelArch != "" {
		meta.KernelArch = hi.KernelArch
	}
	if hi.VirtualizationSystem != "" {
		meta.VirtualizationSystem = hi.VirtualizationSystem
	}
	if hi.VirtualizationRole != "" {
		meta.VirtualizationRole = hi.VirtualizationRole
	}

	return meta, nil

}
