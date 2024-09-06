package agent

import (
	"big-john/internal/ai"
	db "big-john/internal/db/postgresql/sqlc"
	"big-john/pkg/logger"
	"context"
	"encoding/json"
	"strings"
)

type Agent interface {
	ProcessInput(input string) (string, error)
}

type BaseAgent struct {
	aiAdapter *ai.Adapter
	store     db.Store
}

func (a *BaseAgent) ProcessInput(input string) (string, error) {
	log := logger.Get()
	log.Info().Str("input", input).Msg("Processing input")

	response, err := a.aiAdapter.ProcessPrompt(input)

	if err != nil {
		log.Error().Err(err).Msg("Error processing with AI")
		return "", err
	}

	// Store the result - REDIS

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

func NewAgent(aiModel *ai.Adapter, store db.Store) *BaseAgent {
	return &BaseAgent{
		aiAdapter:  aiModel,
		store: store,
	}
}

type CategoryAgent struct {
	BaseAgent
	categories []string
}

func NewCategoryAgent(aiAdapter *ai.Adapter, store db.Store, categories []string) *CategoryAgent {
	return &CategoryAgent{
		BaseAgent:  *NewAgent(aiAdapter, store),
		categories: categories,
	}
}

func (a *CategoryAgent) ProcessInput(input string) (string, error) {
	log := logger.Get()
	log.Info().Str("input", input).Msg("Categorizing prompt")
	service_providers, err := a.store.ListServiceProviders(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Database query failed")
		return "", err
	}
	
	for _, service_provider := range service_providers {
		log.Info().
			Str("name", service_provider.Name).
			Str("speciality", service_provider.Specialty.String).
			Msg("Service Provider details")
	}

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
