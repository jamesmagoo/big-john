package api

import (
	"big-john/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func pingHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := r.PathValue("name")
	w.Write([]byte("pong" + "_" + pathParams))
}

func (s *APIServer) handlePrompt(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Prompt string `json:"prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := s.processor.ProcessPrompt(req.Prompt)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to process prompt")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	log.Info().Str("URL", r.URL.Path)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func (s *APIServer) handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		s.log.Error().Err(err).Msg("Failed to parse Telegram update")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if update.Message == nil {
		return
	}

	response, err := s.processor.ProcessPrompt(update.Message.Text)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to process Telegram message")
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	_, err = s.telegramBot.Send(msg)
	
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to send response Telegram message")
	}
}