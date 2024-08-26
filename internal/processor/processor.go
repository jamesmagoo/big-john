package processor

import (
	"big-john/internal/agent"
	"big-john/pkg/logger"
)

type Processor struct {
	agent *agent.Agent
	log   *logger.Logger
}

func NewProcessor(a *agent.Agent, l *logger.Logger) *Processor {
	return &Processor{
		agent: a,
		log:   l,
	}
}

func (p *Processor) ProcessPrompt(prompt string) (string, error) {
	// For now, we're just passing the prompt directly to the agent
	// In the future, we can add pre-processing and post-processing logic here
	p.log.Info().Str("prompt", prompt).Msg("Processing prompt")
	return p.agent.ProcessInput(prompt)
}