package api

import (
	"big-john/internal/processor"
	"big-john/pkg/logger"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type APIServer struct {
	addr      string
	processor *processor.Processor
	log       *logger.Logger
	hub       *Hub
	telegramBot      *tgbotapi.BotAPI
}

func NewAPIServer(addr string, p *processor.Processor) *APIServer {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_AUTH_TOKEN"))
	if err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to create Telegram bot")
	}
	return &APIServer{
		addr:      addr,
		processor: p,
		log:       logger.Get(),
		hub:       newHub(),
		telegramBot: bot,
	}
}

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
		s.log.Error().Err(err).Msg("Failed to send Telegram message")
	}
}

func (s *APIServer) Run() error {

	go s.hub.run()

	router := http.NewServeMux()

	router.HandleFunc("/", serveHome)
	router.HandleFunc("POST /prompt", s.handlePrompt)
	router.HandleFunc("POST /telegram/webhook", s.handleTelegramWebhook)

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.serveWs(s.hub, w, r)
	})
	router.HandleFunc("GET /users/{uid}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("uid")
		w.Write([]byte("User ID:" + userID))
	})
	router.HandleFunc("POST /ping/{name}", pingHandler)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// order matters here...
	middlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		//RequireAuthMiddleware,
	)

	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(v1),
	}

	logger.PrintAsciiArt()

	s.log.Info().Str("port", server.Addr).Msg("BIG JOHN serving...")

	return server.ListenAndServe()
}
