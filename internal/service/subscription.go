package service

import (
	"Test_EM/internal/models"
	"Test_EM/internal/repository"
	"Test_EM/internal/repository/postgres"
	"context"
	"fmt"
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
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required for delete")
	}
	return s.repo.Delete(ctx, id)
}
