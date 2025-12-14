package main

import (
	"Currency-apiNew2/internal/currency/transport/http"
	"Currency-apiNew2/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	log := logger.New()

	srv := http.NewServer(log)
	if err := srv.Run(); err != nil {
		log.Fatal("currency init failed", zap.Error(err))
	}
}
