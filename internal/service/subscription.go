package service

import (
	"Test_EM/internal/models"
	"Test_EM/internal/repository"
	"Test_EM/internal/repository/postgres"
	"context"
	"fmt"
	"time"
)

type SubscriptionService struct {
	repo repository.Subscription
}

func NewSubscriptionService(r *postgres.Repository) *SubscriptionService {
	return &SubscriptionService{repo: r}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub models.Subscription) error {
	if sub.Price < 0 {
		return fmt.Errorf("цена не может быть отрицательной")
	}
	if sub.ServiceName == "" {
		return fmt.Errorf("название сервиса обязательно")
	}
	if sub.UserID == "" {
		return fmt.Errorf("ID пользователя обязателен")
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return fmt.Errorf("end_date cannot be earlier than start_date")
	}

	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	return s.repo.GetAll(ctx)
}

func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id string) (models.Subscription, error) {
	if id == "" {
		return models.Subscription{}, fmt.Errorf("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, sub models.Subscription) error {
	if sub.ID == "" {
		return fmt.Errorf("id is required for update")
	}

	if sub.Price < 0 {
		return fmt.Errorf("the price cannot be negative")
	}
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return fmt.Errorf("end_date cannot be earlier than start_date")
	}
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required for delete")
	}
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) GetStats(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error) {
	if to.Before(from) {
		return 0, fmt.Errorf("to_date cannot be earlier than from_date")
	}

	subs, err := s.repo.GetActiveInPeriod(ctx, userID, serviceName, from, to)
	if err != nil {
		return 0, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	totalCost := 0

	for _, sub := range subs {
		calcStart := sub.StartDate
		if calcStart.Before(from) {
			calcStart = from
		}

		calcEnd := to
		if sub.EndDate != nil && sub.EndDate.Before(to) {
			calcEnd = *sub.EndDate
		}

		monthsActive := (calcEnd.Year()-calcStart.Year())*12 + int(calcEnd.Month()) - int(calcStart.Month()) + 1

		if monthsActive > 0 {
			totalCost += monthsActive * sub.Price
		}
	}

	return totalCost, nil
}
