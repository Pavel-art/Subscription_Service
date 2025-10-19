package core_interfaces

import (
	"SubscriptionService/internal/core/models"
	"SubscriptionService/internal/core/ports/filters"
	"context"

	"github.com/google/uuid"
)

type ISubRepository interface {
	Create(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetAll(ctx context.Context, page, pageSize int64) ([]*models.Subscription, int64, int64, error)
	SumSubscriptionsCost(ctx context.Context, filter *filters.SubFilter) (int64, error)
}
