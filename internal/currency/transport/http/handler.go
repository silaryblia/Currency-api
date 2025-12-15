package http

import (
	"Currency-apiNew2/internal/currency/domain"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	service domain.CurrencyRepository
	logger  *zap.Logger
}

func NewHandler(service domain.CurrencyRepository, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	rates, err := h.service.GetAll()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, rates, "")
}

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := strings.ToUpper(vars["code"])

	rate, err := h.service.GetOne(code)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, nil, "Currency not found")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code": code,
		"rate": rate,
	}, "")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Currency string  `json:"code"`
		Rate     float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid JSON body")
		return
	}

	if req.Currency == "" || req.Rate == 0 {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid parameters")
		return
	}

	if err := h.service.Create(req.Currency, req.Rate); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, req, "")
}

func (h *Handler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Currency string  `json:"code"`
		Rate     float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid JSON")
		return
	}

	if req.Currency == "" {
		WriteJSON(w, http.StatusBadRequest, nil, "currency code required")
		return
	}

	if err := h.service.UpdateOne(req.Currency, req.Rate); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, req, "")
}

func (h *Handler) UpdateAll(w http.ResponseWriter, r *http.Request) {
	h.service.UpdateAll()
	rates, _ := h.service.GetAll()

	WriteJSON(w, http.StatusOK, rates, "")
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	h.service.DeleteAll()

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Все курсы удалены",
	}, "")
}
