package handler

import (
	"encoding/json"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Response универсальный формат ответа (аналог API Resources)
// @name APIResponse
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

type Handler struct {
	Subscription *SubscriptionHandler
}

func NewHandler(sub *SubscriptionHandler) *Handler {
	return &Handler{Subscription: sub}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /subscriptions", h.Subscription.List)
	mux.HandleFunc("POST /subscriptions", h.Subscription.Create)
	mux.HandleFunc("GET /subscriptions/{id}", h.Subscription.Get)
	mux.HandleFunc("PUT /subscriptions/{id}", h.Subscription.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Subscription.Delete)
	mux.HandleFunc("GET /subscriptions/stats", h.Subscription.GetStats)

	return mux
}
