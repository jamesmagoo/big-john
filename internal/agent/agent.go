package agent

import (
	"big-john/internal/ai"
	"big-john/internal/store"
	"big-john/pkg/logger"
	"encoding/json"
	"strings"
)

type Agent interface {
	ProcessInput(input string) (string, error)
}

type BaseAgent struct {
	aiAdapter  *ai.Adapter
	dataSource *data.Source
}


func (a *BaseAgent) ProcessInput(input string) (string, error) {
	log := logger.Get()
	log.Info().Str("input", input).Msg("Processing input")

	// Check if we have a cached response
	cachedResponse, err := a.dataSource.GetCachedResponse(input)
	if err == nil {
		log.Info().Msg("Returning cached response")
		return cachedResponse, nil
	}

	// If no cached response, process with AI
	log.Info().Msg("No cache hit, processing with AI")

	response, err := a.aiAdapter.ProcessPrompt(input)

	if err != nil {
		log.Error().Err(err).Msg("Error processing with AI")
		return "", err
	}

	// Store the result
	err = a.dataSource.StoreResult(input, response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to store result in cache")
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

func NewAgent(aiModel *ai.Adapter, d *data.Source) *BaseAgent {
	return &BaseAgent{
		aiAdapter:  aiModel,
		dataSource: d,
	}
}

type CategoryAgent struct {
	BaseAgent
	categories []string
}

func NewCategoryAgent(aiAdapter *ai.Adapter, dataSource *data.Source, categories []string) *CategoryAgent {
	return &CategoryAgent{
		BaseAgent:  *NewAgent(aiAdapter, dataSource),
		categories: categories,
	}
}

func (a *CategoryAgent) ProcessInput(input string) (string, error) {
	log := logger.Get()
	log.Info().Str("input", input).Msg("Categorizing prompt")

	// Prepare the prompt for the LLM
	prompt := a.prepareCategoryPrompt(input)

	// Process with AI
	response, err := a.aiAdapter.ProcessPrompt(prompt)
	if err != nil {
		log.Error().Err(err).Msg("Error processing with AI")
		return "", err
	}

	// Parse the response
	category, err := a.parseCategory(response)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing category")
		return "", err
	}

	return category, nil
}

func (a *CategoryAgent) prepareCategoryPrompt(input string) string {
	categoriesJSON, _ := json.Marshal(a.categories)
	return `Categorize the following input into one of these categories: ` + string(categoriesJSON) + `
If none of the categories fit, respond with "other".
Input: "` + input + `"
Category:`
}

func (a *CategoryAgent) parseCategory(response string) (string, error) {
	// Simple parsing: just return the trimmed response
	// In a more robust implementation, you might want to validate against the known categories
	return strings.TrimSpace(response), nil
}

// ... (keep other existing code)
