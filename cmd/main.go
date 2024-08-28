package main

import (
	"big-john/internal/agent"
	"big-john/internal/ai"
	"big-john/internal/api"
	"big-john/internal/processor"
	"big-john/internal/store"
	"big-john/pkg/logger"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Load .env file
	log := logger.Get()
	err := godotenv.Load()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal().Msg("OPENAI_API_KEY environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	agentManager := agent.NewAgentManager()

	// Create multiple AI adapters and agents
	aiAdapter1 := ai.NewAdapter("openai", openai.GPT3Dot5Turbo)
	aiAdapter2 := ai.NewAdapter("anthropic", "claude-3.5-sonnet")
	dataSource := data.NewSource()

	agent1 := agent.NewAgent(aiAdapter1, dataSource)
	agent2 := agent.NewAgent(aiAdapter2, dataSource)
	categories := []string{"hair", "nails", "makeup"}

    // Create the CategoryAgent
    categoriserAgent := agent.NewCategoryAgent(aiAdapter1, dataSource, categories)

	agentManager.AddAgent("agent", agent1)
	agentManager.AddAgent("agent2", agent2)
	agentManager.AddAgent("categoriser", categoriserAgent)


	proc := processor.NewProcessor(agentManager)
	server := api.NewAPIServer(":"+port, proc)

	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}