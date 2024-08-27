package main

import (
	"big-john/internal/agent"
	"big-john/internal/ai"
	"big-john/internal/api"
	"big-john/internal/store"
	"big-john/internal/processor"
	"big-john/pkg/logger"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	log := logger.Get()

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal().Msg("OPENAI_API_KEY environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	aiAdapter := ai.NewAdapter()
	dataSource := data.NewSource()
	agent := agent.NewAgent(aiAdapter, dataSource, log)
	proc := processor.NewProcessor(agent, log)
	server := api.NewAPIServer(":"+port, proc)

	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}