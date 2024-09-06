package main

import (
	"big-john/internal/agent"
	"big-john/internal/ai"
	"big-john/internal/api"
	db "big-john/internal/db/postgresql/sqlc"
	"big-john/internal/processor"
	"big-john/internal/util"
	"big-john/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	logger.PrintAsciiArt()
	config, err := util.LoadConfig(".")
	log := logger.Get()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	var pool *pgxpool.Pool

	if config.Env == "local" {
		log.Info().Msg("Using local db source")
		pool, err = pgxpool.New(context.Background(), config.DBSourceLocal)
	} else {
		pool, err = pgxpool.New(context.Background(), config.DBSource)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db")
	}

	err = pool.Ping(context.Background())

	if err != nil {
		log.Error().Err(err).Msg("Pinged db pool fail..")
	} else {
		log.Info().Str("db_source", config.DBSource).Msg("Connected to the DB")
	}

	store := db.NewStore(pool)

	agentManager := agent.NewAgentManager()

	// Create multiple AI adapters and agents
	aiAdapter1 := ai.NewAdapter("openai", openai.GPT3Dot5Turbo, &config)
	aiAdapter2 := ai.NewAdapter("anthropic", "claude-3.5-sonnet", &config)

	agent1 := agent.NewAgent(aiAdapter1, store)
	agent2 := agent.NewAgent(aiAdapter2, store)
	categories := []string{"availability", "service_providers", "appointments"}

	// Create the CategoryAgent
	categoriserAgent := agent.NewCategoryAgent(aiAdapter1, store, categories)

	agentManager.AddAgent("agent", agent1)
	agentManager.AddAgent("agent2", agent2)
	agentManager.AddAgent("categoriser", categoriserAgent)

	proc := processor.NewProcessor(agentManager)
	server := api.NewAPIServer(&config, proc)

	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}
