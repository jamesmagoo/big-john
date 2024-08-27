package api

import (
	"big-john/internal/processor"
	"big-john/pkg/logger"
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
