package processor

import (
	"big-john/internal/agent"
	"big-john/pkg/logger"
)


type Processor struct {
	agentManager *agent.AgentManager
}

func NewProcessor(am *agent.AgentManager) *Processor {
	return &Processor{
		agentManager: am,
	}
}

func (p *Processor) ProcessPrompt(prompt string) (string, error) {
	log := logger.Get()
	log.Info().Str("prompt", prompt).Msg("Processing prompt")

	// Chain agents
	categorizerAgent, err := p.agentManager.GetAgent("categoriser")
	if err != nil {
		return "", err
	}
	baseAgent, err := p.agentManager.GetAgent("agent")
	if err != nil {
		return "", err
	}

	// Call categorizer agent
	category, err := categorizerAgent.ProcessInput(prompt)
	if err != nil {
		return "", err
	}

	log.Info().Str("category", category).Msg("Prompt categorized")

	// Call base agent with categorized prompt
	result, err := baseAgent.ProcessInput(prompt + " [Category: " + category + "]")
	if err != nil {
		return "", err
	}

	return result, nil
}
