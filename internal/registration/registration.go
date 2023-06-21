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
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/circonus/collector-management-agent/internal/config/keys"
	"github.com/circonus/collector-management-agent/internal/credentials"
	"github.com/denisbrodbeck/machineid"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/viper"
)

type Registration struct {
	Meta    Meta    `json:"meta"`
	AWSMeta AWSTags `json:"aws,omitempty"`
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

type AWSTags struct {
	InstanceID       string    `json:"instance_id,omitempty"`
	AccountID        string    `json:"account_id,omitempty"`
	Architecture     string    `json:"architecture,omitempty"`
	AvailabilityZone string    `json:"availability_zone,omitempty"`
	ImageID          string    `json:"image_id,omitempty"`
	InstanceType     string    `json:"instance_type,omitempty"`
	KernelID         string    `json:"kernel_id,omitempty"`
	PendingTime      time.Time `json:"pending_time,omitempty"`
	PrivateIP        string    `json:"private_ip,omitempty"`
	RamdiskID        string    `json:"ramdisk_id,omitempty"`
	Region           string    `json:"region,omitempty"`
	Version          string    `json:"version,omitempty"`
}

type Response struct {
	AuthToken string `json:"auth_token" yaml:"auth_token"`
	Agent     Agent  `json:"agent" yaml:"agent"`
}

type Agent struct {
	ID string `json:"agent_id" yaml:"agent_id"`
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

	reg := Registration{
		Meta: meta,
	}

	awstags := viper.GetStringSlice(keys.AWSEC2Tags)
	if len(awstags) > 0 {
		at, err := getAWSTags(ctx, awstags) //nolint:govet
		if err != nil {
			log.Fatal().Err(err).Msg("adding AWS EC2 tags")
		}
		reg.AWSMeta = at
	}

	jwt, err := getJWT(ctx, token, reg)
	if err != nil {
		log.Fatal().Err(err).Msg("getting token")
	}

	if err := credentials.SaveJWT([]byte(jwt.AuthToken)); err != nil {
		log.Fatal().Err(err).Msg("saving token")
	}
	if err := credentials.SaveAgentID([]byte(jwt.Agent.ID)); err != nil {
		log.Fatal().Err(err).Msg("saving agent id")
	}

	return nil
}

func getJWT(ctx context.Context, token string, reg Registration) (*Response, error) {
	if token == "" {
		return nil, fmt.Errorf("invalid token (empty)") //nolint:goerr113
	}
	if reg.Meta.Hostname == "" {
		return nil, fmt.Errorf("invalid claims (empty hostname)") //nolint:goerr113
	}
	if reg.Meta.MachineID == "" {
		return nil, fmt.Errorf("invalid claims (empty machine id)") //nolint:goerr113
	}

	c, err := json.Marshal(reg)
	if err != nil {
		return nil, fmt.Errorf("marshal claims: %w", err)
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "agent", "register")
	if err != nil {
		return nil, fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(c))
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

func getAWSTags(ctx context.Context, tags []string) (AWSTags, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return AWSTags{}, fmt.Errorf("failed loading default AWS config: %w", err)
	}
	imdsClient := imds.NewFromConfig(cfg)

	iido, err := imdsClient.GetInstanceIdentityDocument(
		ctx,
		&imds.GetInstanceIdentityDocumentInput{},
	)
	if err != nil {
		return AWSTags{}, fmt.Errorf("failed getting instance identity document: %w", err)
	}

	aws := AWSTags{}

	for _, tag := range tags {
		switch tag {
		case "account_id":
			aws.AccountID = iido.AccountID
		case "architecture":
			aws.Architecture = iido.Architecture
		case "availability_zone":
			aws.AvailabilityZone = iido.AvailabilityZone
		case "image_id":
			aws.ImageID = iido.ImageID
		case "instance_id":
			aws.InstanceID = iido.InstanceID
		case "instance_type":
			aws.InstanceType = iido.InstanceType
		case "kernel_id":
			aws.KernelID = iido.KernelID
		case "pending_time":
			aws.PendingTime = iido.PendingTime
		case "private_ip":
			aws.PrivateIP = iido.PrivateIP
		case "ramdisk_id":
			aws.RamdiskID = iido.RamdiskID
		case "region":
			aws.Region = iido.Region
		case "version":
			aws.Version = iido.Version
		}
	}

	return aws, nil
}
