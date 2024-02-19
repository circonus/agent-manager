package registration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/circonus/agent-manager/internal/config/defaults"
	"gopkg.in/yaml.v3"
)

type Agent struct {
	AgentID     string `json:"agent_id"      yaml:"agent_id"`
	AgentTypeID string `json:"agent_type_id" yaml:"agent_type_id"`
}
type Agents []Agent

func LoadInstalledAgents() (Agents, error) {
	agentFile := filepath.Join(defaults.EtcPath, "agents.yaml")

	data, err := os.ReadFile(agentFile)
	if err != nil {
		return nil, err
	}

	var a Agents
	if err := yaml.Unmarshal(data, &a); err != nil {
		return nil, err
	}

	return a, nil
}

func SaveInstalledAgents(a Agents) error {
	agentFile := filepath.Join(defaults.EtcPath, "agents.yaml")

	data, err := yaml.Marshal(a)
	if err != nil {
		return err
	}

	return os.WriteFile(agentFile, data, 0o600)
}

func GetInstalledAgentID(agentTypeID string) (string, error) {
	agents, err := LoadInstalledAgents()
	if err != nil {
		return "", err
	}

	for _, agent := range agents {
		if agent.AgentTypeID == agentTypeID {
			return agent.AgentID, nil
		}
	}

	return "", fmt.Errorf("agent [%s] not found in installed agents", agentTypeID)
}
