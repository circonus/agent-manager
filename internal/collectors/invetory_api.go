package collectors

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type APICollectors []APICollector

type APICollector struct {
	Platforms []Platform `json:"platforms"`
}

type Platform struct {
	CollectorPlatformID string       `json:"collector_platform_id" yaml:"collector_platform_id"`
	CollectorTypeID     string       `json:"collector_type_id" yaml:"collector_type_id"`
	Executable          string       `json:"executable" yaml:"executable"`
	Start               string       `json:"start" yaml:"start"`
	Stop                string       `json:"stop" yaml:"stop"`
	Restart             string       `json:"restart" yaml:"restart"`
	Reload              string       `json:"reload" yaml:"reload"`
	Status              string       `json:"status" yaml:"status"`
	Version             string       `json:"version" yaml:"version"`
	ConfigFiles         []ConfigFile `json:"config_files" yaml:"config_files"`
}

type ConfigFile struct {
	ConfigFileID string `json:"collector_config_file_id" yaml:"collector_config_file_id"`
	Path         string `json:"path" yaml:"path"`
}

func ParseAPICollectors(data []byte) (Collectors, error) {
	var c APICollectors
	if err := json.Unmarshal(data, &c); err != nil {
		log.Error().Err(err).Msg("parsing api collectors json")
	}

	collectors := make(map[string]map[string]Collector)
	for _, collector := range c {
		for _, platform := range collector.Platforms {
			col := Collector{
				Binary:      platform.Executable,
				Start:       platform.Start,
				Stop:        platform.Stop,
				Restart:     platform.Restart,
				Reload:      platform.Reload,
				Version:     platform.Version,
				ConfigFiles: make(map[string]string, len(platform.ConfigFiles)),
			}
			for _, f := range platform.ConfigFiles {
				col.ConfigFiles[f.ConfigFileID] = f.Path
			}
			if _, ok := collectors[platform.CollectorPlatformID]; !ok {
				collectors[platform.CollectorPlatformID] = make(map[string]Collector)
			}
			collectors[platform.CollectorPlatformID][platform.CollectorTypeID] = col
		}
	}

	return collectors, nil
}
