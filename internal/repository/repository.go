package repository

import (
	"Test_EM/internal/models"
	"context"
	"time"
)

type Subscription interface {
	Create(ctx context.Context, sub models.Subscription) (models.Subscription, error)
	GetAll(ctx context.Context) ([]models.Subscription, error)
	GetByID(ctx context.Context, id string) (models.Subscription, error)
	Update(ctx context.Context, sub models.Subscription) error
	Delete(ctx context.Context, id string) error
	GetActiveInPeriod(ctx context.Context, userID string, serviceName string, from, to time.Time) ([]models.Subscription, error)
}
