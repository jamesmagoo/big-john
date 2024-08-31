package main

import (
	"big-john/internal/agent"
	"big-john/internal/ai"
	"big-john/internal/api"
	"big-john/internal/processor"
	data "big-john/internal/store"
	"big-john/internal/util"
	"big-john/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	log := logger.Get()
	config , err := util.LoadConfig(".")

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	_, err = pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db")
	}

	agentManager := agent.NewAgentManager()

	// Create multiple AI adapters and agents
	aiAdapter1 := ai.NewAdapter("openai", openai.GPT3Dot5Turbo, &config)
	aiAdapter2 := ai.NewAdapter("anthropic", "claude-3.5-sonnet", &config)
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
	server := api.NewAPIServer(&config, proc)

	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}