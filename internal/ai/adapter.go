package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
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
}

// NewOpenAIModel creates a new instance of OpenAIModel
func NewOpenAIModel() *OpenAIModel {
	return &OpenAIModel{
		APIKey: os.Getenv("OPENAI_API_KEY"),
	}
}

// ProcessPrompt sends a request to the OpenAI API and returns a structured response
func (o *OpenAIModel) ProcessPrompt(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model":    "gpt-3.5-turbo",
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to call OpenAI API")
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	output := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return output, nil
}