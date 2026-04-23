package handler

import "net/http"

type Handler struct {
	Subscription *SubscriptionHandler
}

func NewHandler(sub *SubscriptionHandler) *Handler {
	return &Handler{Subscription: sub}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /subscriptions", h.Subscription.List)
	mux.HandleFunc("POST /subscriptions", h.Subscription.Create)
	mux.HandleFunc("GET /subscriptions/{id}", h.Subscription.Get)
	mux.HandleFunc("PUT /subscriptions/{id}", h.Subscription.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Subscription.Delete)

	return mux
}
