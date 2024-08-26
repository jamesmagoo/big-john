package agent

import (
	"big-john/internal/ai"
	"big-john/internal/store"
	"big-john/pkg/logger"
)

type Agent struct {
	aiAdapter  *ai.Adapter
	dataSource *data.Source
	log        *logger.Logger
}

func NewAgent(a *ai.Adapter, d *data.Source, l *logger.Logger) *Agent {
	return &Agent{
		aiAdapter:  a,
		dataSource: d,
		log:        l,
	}
}

func (a *Agent) ProcessInput(input string) (string, error) {
	a.log.Info().Str("input", input).Msg("Processing input")

	// Check if we have a cached response
	cachedResponse, err := a.dataSource.GetCachedResponse(input)
	if err == nil {
		a.log.Info().Msg("Returning cached response")
		return cachedResponse, nil
	}

	// If no cached response, process with AI
	a.log.Info().Msg("No cache hit, processing with AI")
	response, err := a.aiAdapter.ProcessPrompt(input)
	if err != nil {
		a.log.Error().Err(err).Msg("Error processing with AI")
		return "", err
	}

	// Store the result
	err = a.dataSource.StoreResult(input, response)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to store result in cache")
		// We don't return this error because the AI processing was successful
	}

	return response, nil
}
