package handler

import (
	"Test_EM/internal/models"
	"Test_EM/internal/service"
	"encoding/json"
	"net/http"
	"time"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(s *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: s}
}

const dateLayout = "01-2006"

// Create godoc
// @Summary      Создать подписку
// @Description  Создает новую запись о подписке пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body models.CreateSubscriptionRequest true "Данные подписки"
// @Success      201  {object}  map[string]string
// @Failure      400  {string}  string "invalid request body"
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	start, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format, use MM-YYYY", http.StatusBadRequest)
		return
	}

	sub := models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
	}

	if req.EndDate != "" {
		end, err := time.Parse(dateLayout, req.EndDate)
		if err != nil {
			http.Error(w, "invalid end_date format", http.StatusBadRequest)
			return
		}
		sub.EndDate = &end
	}

	if err := h.service.CreateSubscription(r.Context(), sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})

}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := h.service.GetAllSubscriptions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

// Get — получение одной записи по ID
func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Достаем {id} из маршрута (Go 1.22)
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// Update — обновление записи
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.UpdateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Парсим даты, если они пришли (в реальном коде лучше вынести в хелпер)
	start, _ := time.Parse("01-2006", req.StartDate)

	sub := models.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   start,
	}

	if err := h.service.UpdateSubscription(r.Context(), sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}

// Delete — удаление записи
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteSubscription(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
