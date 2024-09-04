package ai

import (
	db "big-john/internal/db/postgresql/sqlc"
	"big-john/internal/util"
	"big-john/pkg/logger"
	"context"

	openai "github.com/sashabaranov/go-openai"
)

// AIModel defines the interface that all AI models must implement
type AIModel interface {
	ProcessPrompt(prompt string) (string, error)
}

// Adapter wraps an AIModel and provides additional functionality
type Adapter struct {
	aiServiceProvider AIModel
	modelName         string
	config            *util.Config
	store             db.Store
}

// NewAdapter creates a new instance of Adapter with an AIModel
func NewAdapter(modelType string, modelName string, config *util.Config, store db.Store) *Adapter {
	var model AIModel
	switch modelType {
	case "openai":
		model = NewOpenAIModel(modelName, config, store)
	case "anthropic":
		model = NewAnthropicModel(modelName, config, store)
	default:
		model = NewOpenAIModel(modelName, config, store)
	}
	return &Adapter{
		aiServiceProvider: model,
		modelName:         modelName,
		config:            config,
		store:             store,
	}
}

// ProcessPrompt sends a prompt to the AI model and returns the response
func (a *Adapter) ProcessPrompt(prompt string) (string, error) {
	return a.aiServiceProvider.ProcessPrompt(prompt)
}

// OpenAIModel is an implementation of the AIModel interface for OpenAI
type OpenAIModel struct {
	APIKey    string
	log       *logger.Logger
	modelName string
	config    *util.Config
	store     db.Store
}

// NewOpenAIModel creates a new instance of OpenAIModel
func NewOpenAIModel(modelName string, config *util.Config, store db.Store) *OpenAIModel {
	return &OpenAIModel{
		APIKey:    config.OpenAIAPIKey,
		log:       logger.Get(),
		modelName: modelName,
		config:    config,
		store:     store,
	}
}

// ProcessPrompt sends a request to the OpenAI API and returns a structured response
func (o *OpenAIModel) ProcessPrompt(prompt string) (string, error) {
	// Example database query
	authors, err := o.store.ListServiceProviders(context.Background())
	if err != nil {
		o.log.Error().Err(err).Msg("Database query failed")
		return "", err
	}
	
	// Range over authors and print
	for _, author := range authors {
		o.log.Info().
			Str("name", author.Name).
			Msg("Author details")
	}

	client := openai.NewClient(o.APIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: o.modelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		o.log.Error().Err(err).Msg("Chat completion problem")
		return "", err
	}

	output := resp.Choices[0].Message.Content

	return output, nil
}

// Add a new AI model implementation, e.g., AnthropicModel
type AnthropicModel struct {
	APIKey    string
	log       *logger.Logger
	modelName string
	config    *util.Config
	store     db.Store
}

// NewAnthropicModel creates a new instance of AnthropicModel
func NewAnthropicModel(modelName string, config *util.Config, store db.Store) *AnthropicModel {
	return &AnthropicModel{
		APIKey:    config.OpenAIAPIKey, // TODO specific api key needed...
		log:       logger.Get(),
		modelName: modelName,
		config:    config,
		store:     store,
	}
}

// ProcessPrompt sends a request to the Anthropic API and returns a structured response
func (a *AnthropicModel) ProcessPrompt(prompt string) (string, error) {
	// Implement Anthropic API call here
	// ...
	return "anthropic impl TODO", nil
}

// ... existing code ...
