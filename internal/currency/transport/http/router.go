package http

import (
	"Currency-apiNew2/internal/currency/domain"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(svc domain.CurrencyRepository, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply middlewares
	r.Use(JSONMiddleware)
	r.Use(recoveryMiddleware(logger))

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
