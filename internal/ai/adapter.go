package ai

import (
	"big-john/pkg/logger"
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// AIModel defines the interface that all AI models must implement
type AIModel interface {
	ProcessPrompt(prompt string) (string, error)
}

// Adapter wraps an AIModel and provides additional functionality
type Adapter struct {
	model AIModel
}

// NewAdapter creates a new instance of Adapter with an OpenAIModel
func NewAdapter() *Adapter {
	return &Adapter{
		model: NewOpenAIModel(),
	}
}

// ProcessPrompt sends a prompt to the AI model and returns the response
func (a *Adapter) ProcessPrompt(prompt string) (string, error) {
	return a.model.ProcessPrompt(prompt)
}

// OpenAIModel is an implementation of the AIModel interface for OpenAI
type OpenAIModel struct {
	APIKey string
	log       *logger.Logger
}

// NewOpenAIModel creates a new instance of OpenAIModel
func NewOpenAIModel() *OpenAIModel {
	return &OpenAIModel{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		log: logger.Get(),
	}
}

// ProcessPrompt sends a request to the OpenAI API and returns a structured response
func (o *OpenAIModel) ProcessPrompt(prompt string) (string, error) {
	
	client := openai.NewClient(o.APIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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
    APIKey string
    log    *logger.Logger
}

// NewAnthropicModel creates a new instance of AnthropicModel
func NewAnthropicModel() *AnthropicModel {
    return &AnthropicModel{
        APIKey: os.Getenv("ANTHROPIC_API_KEY"),
        log:    logger.Get(),
    }
}

// ProcessPrompt sends a request to the Anthropic API and returns a structured response
func (a *AnthropicModel) ProcessPrompt(prompt string) (string, error) {
    // Implement Anthropic API call here
    // ...
	return "anthropic impl TODO", nil
}

// ... existing code ...