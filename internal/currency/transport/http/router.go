package http

import (
	"Currency-apiNew2/internal/currency/repository/memory"
	"Currency-apiNew2/internal/currency/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply middlewares
	r.Use(JSONMiddleware)
	r.Use(recoveryMiddleware(logger))

	repo := memory.NewCurrencyRepoInMemory(logger)
	svc := service.NewCurrencyService(repo)
	h := NewHandler(svc, logger)

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/currency", h.GetAll).Methods("GET")
	api.HandleFunc("/currency/{code}", h.GetOne).Methods("GET")
	api.HandleFunc("/currency/create", h.Create).Methods("POST")
	api.HandleFunc("/currency/update", h.UpdateOne).Methods("PUT")
	api.HandleFunc("/update", h.UpdateAll).Methods("PUT")
	api.HandleFunc("/delete", h.DeleteAll).Methods("DELETE")

	return r
}
