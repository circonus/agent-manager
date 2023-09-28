// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

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
	"unicode"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/circonus/agent-manager/internal/config/keys"
	"github.com/circonus/agent-manager/internal/credentials"
	"github.com/circonus/agent-manager/internal/release"
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/spf13/viper"
)

type Registration struct {
	Version              string   `json:"version"`
	MachineID            string   `json:"machine_id"`
	Hostname             string   `json:"hostname"`
	OS                   string   `json:"os"`
	Platform             string   `json:"platform"`
	PlatformVersion      string   `json:"platform_version"`
	PlatformFamily       string   `json:"platform_family"`
	KernelArch           string   `json:"kernel_arch"`
	KernelVersion        string   `json:"kernel_version"`
	VirtualizationSystem string   `json:"virtualization_system"`
	VirtualizationRole   string   `json:"virtualization_role"`
	Data                 Data     `json:"data,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
}

type Data struct {
	AWSMeta AWSTags `json:"aws,omitempty"`
}

type AWSTags struct {
	InstanceID       string `json:"instance_id,omitempty"`
	AccountID        string `json:"account_id,omitempty"`
	Architecture     string `json:"architecture,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
	ImageID          string `json:"image_id,omitempty"`
	InstanceType     string `json:"instance_type,omitempty"`
	KernelID         string `json:"kernel_id,omitempty"`
	PendingTime      string `json:"pending_time,omitempty"`
	PrivateIP        string `json:"private_ip,omitempty"`
	RamdiskID        string `json:"ramdisk_id,omitempty"`
	Region           string `json:"region,omitempty"`
	Version          string `json:"version,omitempty"`
}

type Response struct {
	AuthToken    string `json:"access_token"  yaml:"access_token"`
	ManagerID    string `json:"manager_id"    yaml:"manager_id"`
	RefreshToken string `json:"refresh_token" yaml:"refresh_token"`
}

const (
	maxTags   = 32
	maxTagLen = 256
)

// Start the registration process.
func Start(ctx context.Context) error {
	log.Info().Msg("starting registration")

	token := viper.GetString(keys.Register)

	hn, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("getting hostname")
	}

	if viper.GetString(keys.InstanceID) != "" {
		hn = viper.GetString(keys.InstanceID)
	}

	if hn == "" {
		log.Fatal().Str("hostname", hn).Msg("empty hostname")
	}

	mid, err := getMachineID()
	if err != nil {
		log.Fatal().Err(err).Str("mid", mid).Msg("invalid machine id")
	}

	reg, err := getHostInfo()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to retrieve host info")
	}

	reg.Hostname = hn
	reg.MachineID = mid
	reg.Version = "v" + release.VERSION

	awstags := viper.GetStringSlice(keys.AWSEC2Tags)
	if len(awstags) > 0 {
		at, err := getAWSTags(ctx, awstags) //nolint:govet
		if err != nil {
			log.Fatal().Err(err).Msg("adding AWS EC2 tags")
		}

		reg.Data.AWSMeta = at
	}

	tags := viper.GetStringSlice(keys.Tags)
	if len(tags) > 0 {
		reg.Tags = formatTags(tags)
	}

	jwt, err := getJWT(ctx, token, reg)
	if err != nil {
		log.Fatal().Err(err).Msg("getting token")
	}

	if err := credentials.SaveJWT([]byte(jwt.AuthToken)); err != nil {
		log.Fatal().Err(err).Msg("saving token")
	}

	if err := credentials.SaveRefreshToken([]byte(jwt.RefreshToken)); err != nil {
		log.Fatal().Err(err).Msg("saving token")
	}

	if err := credentials.SaveManagerID([]byte(jwt.ManagerID)); err != nil {
		log.Fatal().Err(err).Msg("saving manager id")
	}

	return nil
}

func getJWT(ctx context.Context, token string, reg Registration) (*Response, error) {
	if token == "" {
		return nil, fmt.Errorf("invalid token (empty)") //nolint:goerr113
	}

	if reg.Hostname == "" {
		return nil, fmt.Errorf("invalid claims (empty hostname)") //nolint:goerr113
	}

	if reg.MachineID == "" {
		return nil, fmt.Errorf("invalid claims (empty machine id)") //nolint:goerr113
	}

	c, err := json.Marshal(reg)
	if err != nil {
		return nil, fmt.Errorf("marshal claims: %w", err)
	}

	reqURL, err := url.JoinPath(viper.GetString(keys.APIURL), "manager", "register")
	if err != nil {
		return nil, fmt.Errorf("req url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(c))
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

func getMachineID() (string, error) {
	if viper.GetBool(keys.UseMachineID) {
		id, err := machineid.ID()
		if err != nil {
			return "", err //nolint:wrapcheck
		}

		mac := hmac.New(sha256.New, []byte(id))

		return fmt.Sprintf("%x", mac.Sum(nil)), nil
	}

	if err := credentials.LoadMachineID(); err != nil {
		// if it doesn't exist, we want to create it
		// otherwise return the actual error
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("loading machine id: %w", err)
		}
	} else if viper.GetString(keys.MachineID) != "" {
		return viper.GetString(keys.MachineID), nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generating uuid: %w", err)
	}

	ids := id.String()

	if err := credentials.SaveMachineID([]byte(ids)); err != nil {
		return "", fmt.Errorf("saving machine id: %w", err)
	}

	return ids, nil
}

func getHostInfo() (Registration, error) {
	reg := Registration{}

	hi, err := host.Info()
	if err != nil {
		return reg, err //nolint:wrapcheck
	}

	// hi.OS is runtime.GOOS
	if hi.OS != "" {
		reg.OS = hi.OS
	}

	if hi.Platform != "" {
		reg.Platform = hi.Platform
	}

	if hi.PlatformFamily != "" {
		reg.PlatformFamily = hi.PlatformFamily
	}

	if hi.PlatformVersion != "" {
		reg.PlatformVersion = hi.PlatformVersion
	}

	if hi.KernelVersion != "" {
		reg.KernelVersion = hi.KernelVersion
	}

	if hi.KernelArch != "" {
		reg.KernelArch = hi.KernelArch
	}

	if hi.VirtualizationSystem != "" {
		reg.VirtualizationSystem = hi.VirtualizationSystem
	}

	if hi.VirtualizationRole != "" {
		reg.VirtualizationRole = hi.VirtualizationRole
	}

	return reg, nil
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
			aws.PendingTime = iido.PendingTime.String()
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

func formatTags(tags []string) []string {
	if len(tags) > maxTags {
		log.Fatal().Int("max_tags", maxTags).Int("num_tags", len(tags)).Msg("too many tags")
	}

	t := make([]string, 0, len(tags))

	for _, tag := range tags {
		if len(tag) > maxTagLen {
			log.Warn().Int("max_tag_len", maxTagLen).Int("tag_len", len(tag)).Str("tag", tag).Msg("tag too long, ignoring")

			continue
		}

		tagOK := true

		for i := 0; i < len(tag); i++ {
			if !unicode.IsPrint(rune(tag[i])) {
				log.Warn().Int("pos", i).Str("tag", tag).Msg("tag contains non-printable character, ignoring")

				tagOK = false

				break
			}
		}

		if tagOK {
			t = append(t, tag)
		}
	}

	return t
}
