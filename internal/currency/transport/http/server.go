package http

import (
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/repository"
	_ "Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	service *service.CurrencyService
	router  *mux.Router
	logger  *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	var repo domain.CurrencyRepository

	if os.Getenv("USE_POSTGRES") == "true" {
		db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
		if err != nil {
			logger.Fatal("failed to connect postgres", zap.Error(err))
		}

		repo = repository.NewCurrencyRepoPostgres(db, logger)
		logger.Info("using Postgres repository")
	} else {
		repo = repository.NewCurrencyRepoInMemory(logger)
		logger.Info("using InMemory repository")
	}

	svc := service.NewCurrencyService(repo)
	r := NewRouter(svc, logger)

	return &Server{
		router: r,
		logger: logger,
	}
}

func (s *Server) Run() error {
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	go func() {
		s.logger.Info("Server starting...", zap.String("addr", httpServer.Addr))

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server run error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(ctx)
}
