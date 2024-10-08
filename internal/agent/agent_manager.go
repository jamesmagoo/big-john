package agent

import (
	"errors"
)

type AgentManager struct {
	agents map[string]Agent 
}

func NewAgentManager() *AgentManager {
	return &AgentManager{
		agents: make(map[string]Agent), 
	}
}

func (am *AgentManager) AddAgent(id string, agent Agent) {
	am.agents[id] = agent 
}

func (am *AgentManager) GetAgent(id string) (Agent, error) {
	agent, exists := am.agents[id]
	if !exists {
		return nil, errors.New("agent not found")
	}
	return agent, nil
}

func (am *AgentManager) GetAllAgents() []Agent {
	agents := make([]Agent, 0, len(am.agents))
	for _, agent := range am.agents {
		agents = append(agents, agent)
	}
	return agents
}