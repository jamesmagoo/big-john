package api

import (
	"big-john/internal/processor"
	"big-john/internal/util"
	"big-john/pkg/logger"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type APIServer struct {
	config      *util.Config
	processor   *processor.Processor
	log         *logger.Logger
	hub         *Hub
	telegramBot *tgbotapi.BotAPI
}

func NewAPIServer(config *util.Config, p *processor.Processor) *APIServer {
	// TODO should probably not make Telegram bot here? not sure....
	bot, err := tgbotapi.NewBotAPI(config.TelegramAuthToken)
	if err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to create Telegram bot")
	}
	return &APIServer{
		config:      config,
		processor:   p,
		log:         logger.Get(),
		hub:         newHub(),
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
		Addr:    ":"+s.config.ServerAddress,
		Handler: middlewareChain(v1),
	}

	s.log.Info().Str("address", server.Addr).Msg("BIG JOHN serving...")

	return server.ListenAndServe()
}
