package app_interfaces

import (
	"SubscriptionService/internal/api/dto"
	"SubscriptionService/internal/core/models"
	"context"

	"github.com/google/uuid"
)

type ISubService interface {
	Create(ctx context.Context, req dto.CreateSubscriptionRequest) (*models.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateSubscriptionRequest) (*models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetAll(ctx context.Context, page, pageSize int64) (dto.GetAllResponse, error)
	CalculateTotalCost(ctx context.Context, req dto.CostCalculationQueryRequest) (int64, error)
}
