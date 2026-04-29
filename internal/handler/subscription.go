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
// @Success      201  {object}  Response{data=map[string]string} "Успешное создание"
// @Failure      400  {object}  Response "Ошибка в запросе"
// @Failure      500  {object}  Response "Ошибка сервера"
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubscriptionRequest

	// Используем Decoder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Error: "invalid request body"})
		return
	}

	// Валидация даты (можно потом вынести в сервис или хелпер)
	start, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Error: "invalid start_date format, use MM-YYYY"})
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
			sendJSON(w, http.StatusBadRequest, Response{Error: "invalid end_date format"})
			return
		}
		sub.EndDate = &end
	}

	// Вызов сервиса
	if err := h.service.CreateSubscription(r.Context(), sub); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	// Красивый ответ в стиле Resource
	sendJSON(w, http.StatusCreated, Response{
		Message: "Subscription created successfully",
		Data:    map[string]string{"status": "created"},
	})
}

// List godoc
// @Summary      Получить список подписок
// @Description  Возвращает все существующие подписки
// @Tags         subscriptions
// @Produce      json
// @Success      200  {object}  Response{data=[]models.Subscription} "Список подписок"
// @Failure      500  {object}  Response "Ошибка сервера"
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	subs, err := h.service.GetAllSubscriptions(r.Context())
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Data: subs,
	})
}

// Get godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает одну запись о подписке
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки (UUID)"
// @Success      200  {object}  Response{data=models.Subscription} "Данные подписки"
// @Failure      404  {object}  Response "Подписка не найдена"
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		sendJSON(w, http.StatusBadRequest, Response{Error: "missing id"})
		return
	}

	sub, err := h.service.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		sendJSON(w, http.StatusNotFound, Response{Error: "subscription not found"})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Data: sub,
	})
}

// Update godoc
// @Summary      Обновить подписку
// @Description  Обновляет данные существующей подписки
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "ID подписки (UUID)"
// @Param        input body      models.UpdateSubscriptionRequest  true  "Новые данные"
// @Success      200  {object}  Response{data=models.Subscription} "Обновленная подписка"
// @Failure      400  {object}  Response "Ошибка валидации"
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.UpdateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Error: "bad request body"})
		return
	}

	currentSub, err := h.service.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		sendJSON(w, http.StatusNotFound, Response{Error: "subscription not found"})
		return
	}

	if req.ServiceName != nil {
		currentSub.ServiceName = *req.ServiceName
	}

	if req.Price != nil {
		currentSub.Price = *req.Price
	}

	if req.StartDate != nil {
		start, err := time.Parse(dateLayout, *req.StartDate)
		if err != nil {
			sendJSON(w, http.StatusBadRequest, Response{Error: "invalid start_date format"})
			return
		}
		currentSub.StartDate = start
	}

	if req.EndDate != nil {
		if *req.EndDate == "" {
			currentSub.EndDate = nil
		} else {
			end, err := time.Parse(dateLayout, *req.EndDate)
			if err != nil {
				sendJSON(w, http.StatusBadRequest, Response{Error: "invalid end_date format"})
				return
			}
			currentSub.EndDate = &end
		}
	}

	if err := h.service.UpdateSubscription(r.Context(), currentSub); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Message: "Subscription updated successfully",
		Data:    currentSub,
	})
}

// Delete godoc
// @Summary      Удалить подписку
// @Description  Удаляет запись из базы по ID
// @Tags         subscriptions
// @Param        id   path      string  true  "ID подписки (UUID)"
// @Success      200  {object}  Response "Успешное удаление"
// @Failure      500  {object}  Response "Ошибка удаления"
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteSubscription(r.Context(), id); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Message: "Subscription deleted successfully",
	})
}

// GetStats godoc
// @Summary      Получить статистику затрат
// @Description  Считает суммарную стоимость подписок за период с фильтрами
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  true   "ID пользователя"
// @Param        service_name  query     string  false  "Название сервиса"
// @Param        from          query     string  true   "Начало периода (MM-YYYY)"
// @Param        to            query     string  true   "Конец периода (MM-YYYY)"
// @Success      200  {object}  Response{data=map[string]int}
// @Router       /subscriptions/stats [get]
func (h *SubscriptionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userID := query.Get("user_id")
	serviceName := query.Get("service_name")
	from := query.Get("from")
	to := query.Get("to")

	// Базовая валидация обязательных полей
	if userID == "" || from == "" || to == "" {
		sendJSON(w, http.StatusBadRequest, Response{Error: "user_id, from and to are required parameters"})
		return
	}

	total, err := h.service.GetStats(r.Context(), userID, serviceName, from, to)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Data: map[string]int{"total_cost": total},
	})
}
