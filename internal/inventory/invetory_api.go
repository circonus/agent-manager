package inventory

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type APIAgents []APIAgent

type APIAgent struct {
	Platforms []Platform `json:"platforms"`
}

type Platform struct {
	ID          string       `json:"platform_id"   yaml:"platform_id"`
	AgentTypeID string       `json:"agent_type_id" yaml:"agent_type_id"`
	Executable  string       `json:"executable"    yaml:"executable"`
	Commands    []Commands   `json:"commands"      yaml:"commands"`
	Start       string       `json:"start"         yaml:"start"`
	Stop        string       `json:"stop"          yaml:"stop"`
	Restart     string       `json:"restart"       yaml:"restart"`
	Reload      string       `json:"reload"        yaml:"reload"`
	Status      string       `json:"status"        yaml:"status"`
	Version     string       `json:"version"       yaml:"version"`
	ConfigFiles []ConfigFile `json:"config_files"  yaml:"config_files"`
}

type Commands struct {
	Command string `json:"command" yaml:"command"`
	Name    string `json:"name"    yaml:"name"`
}

type ConfigFile struct {
	ConfigFileID string `json:"config_file_id" yaml:"config_file_id"`
	Path         string `json:"path"           yaml:"path"`
}

func ParseAPIAgents(data []byte) (Agents, error) {
	var aa APIAgents
	if err := json.Unmarshal(data, &aa); err != nil {
		log.Error().Err(err).Msg("parsing api agents json")
	}

	agents := make(map[string]map[string]Agent)

	for _, agent := range aa {
		for _, platform := range agent.Platforms {
			col := Agent{
				Binary:      platform.Executable,
				ConfigFiles: make(map[string]string, len(platform.ConfigFiles)),
			}

			for _, c := range platform.Commands {
				switch c.Name {
				case "start":
					col.Start = c.Command
				case "stop":
					col.Stop = c.Command
				case "restart":
					col.Restart = c.Command
				case "reload":
					col.Reload = c.Command
				case "status":
					col.Status = c.Command
				case "version":
					col.Version = c.Command
				default:
					log.Warn().Str("cmd", c.Name).Msg("unknown command")
				}
			}

			for _, f := range platform.ConfigFiles {
				col.ConfigFiles[f.ConfigFileID] = f.Path
			}

			if _, ok := agents[platform.ID]; !ok {
				agents[platform.ID] = make(map[string]Agent)
			}

			agents[platform.ID][platform.AgentTypeID] = col
		}
	}

	return agents, nil
}
