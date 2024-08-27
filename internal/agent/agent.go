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

func NewAgent(aiModel *ai.Adapter, d *data.Source, l *logger.Logger) *Agent {
	return &Agent{
		aiAdapter:  aiModel,
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

	// TODO: Implement additional AI agent capabilities:
	// 1. Natural language understanding to categorize user intents
	// 2. Context management to maintain conversation history
	// 3. Multi-turn dialogue handling for more interactive conversations
	// 4. Integration with external APIs for real-time data retrieval
	// 5. Sentiment analysis to gauge user emotions and adjust responses
	// 6. Personalization based on user preferences or past interactions
	// 7. Multilingual support for global user base
	// 8. Task completion tracking for complex, multi-step requests
	// 9. Proactive suggestions or follow-up questions to enhance user experience
	// 10. Continuous learning from user interactions to improve responses over time

	return response, nil
}
