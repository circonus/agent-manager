// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

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
		return nil, err //nolint:wrapcheck
	}

	var a Agents
	if err := yaml.Unmarshal(data, &a); err != nil {
		return nil, err //nolint:wrapcheck
	}

	return a, nil
}

func SaveInstalledAgents(a Agents) error {
	agentFile := filepath.Join(defaults.EtcPath, "agents.yaml")

	data, err := yaml.Marshal(a)
	if err != nil {
		return err //nolint:wrapcheck
	}

	return os.WriteFile(agentFile, data, 0600) //nolint:wrapcheck
}

func GetInstalledAgentID(agentTypeID string) (string, error) {
	agents, err := LoadInstalledAgents()
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	for _, agent := range agents {
		if agent.AgentTypeID == agentTypeID {
			return agent.AgentID, nil
		}
	}

	return "", fmt.Errorf("agent [%s] not found in installed agents", agentTypeID) //nolint:goerr113
}