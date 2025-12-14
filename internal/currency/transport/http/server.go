package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	router *mux.Router
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	s := &Server{
		router: NewRouter(logger),
		logger: logger,
	}
	return s
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
