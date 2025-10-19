package services

import (
	"SubscriptionService/internal/api/dto"
	appInterfaces "SubscriptionService/internal/application/app_interfaces"
	"SubscriptionService/internal/core/core_interfaces"
	"SubscriptionService/internal/core/models"
	"SubscriptionService/internal/core/ports/filters"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SubService struct {
	repo   core_interfaces.ISubRepository
	logger *zerolog.Logger
}

var _ appInterfaces.ISubService = (*SubService)(nil)

func NewSubService(repo core_interfaces.ISubRepository, logger *zerolog.Logger) *SubService {
	return &SubService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SubService) Create(ctx context.Context, req dto.CreateSubscriptionRequest) (*models.Subscription, error) {
	s.logger.Debug().
		Str("serviceName", req.ServiceName).
		Str("userId", req.UserID.String()).
		Msg("Creating subscription")

	sub, err := models.NewSubscription(
		req.ServiceName,
		req.Price,
		req.UserID,
		req.StartDate,
		req.EndDate,
	)

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("service", req.ServiceName).
			Msg("Failed to create subscription model")
		return nil, err
	}

	createdSub, err := s.repo.Create(ctx, sub)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("subscriptionId", sub.Id.String()).
			Msg("Failed to save subscription to repository")
		return nil, err
	}
	s.logger.Info().
		Str("subscriptionId", createdSub.Id.String()).
		Str("userId", createdSub.UserId.String()).
		Int64("price", createdSub.Price).
		Msg("Subscription created successfully")

	return createdSub, nil
}

func (s *SubService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateSubscriptionRequest) (*models.Subscription, error) {

	s.logger.Debug().
		Str("id", id.String()).
		Msg("Update subscription: started")

	// 1. Получаем существующую запись
	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("id", id.String()).
			Msg("Update subscription: failed to get existing")

		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	if existing == nil {
		s.logger.Warn().
			Str("id", id.String()).
			Msg("Update subscription: not found")
		return nil, fmt.Errorf("Update subscription: not found id ")
	}

	// 2. Применяем partial update
	updated := s.applyPartialUpdate(existing, req)
	s.logger.Debug().
		Str("id", id.String()).
		Msg("Update subscription: partial update applied")

	// 3. Валидируем модель
	if err := updated.Validate(); err != nil {
		s.logger.Warn().
			Err(err).
			Str("id", id.String()).
			Msg("Update subscription: validation failed")
		return nil, err
	}

	// 4. Сохраняем в репозитории
	result, err := s.repo.Update(ctx, updated)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("id", id.String()).
			Msg("Update subscription: failed to update subscription")
		return nil, fmt.Errorf("repository error: %w", err)
	}

	s.logger.Info().
		Str("id", id.String()).
		Msg("Update subscription: success")
	return result, nil
}

func (s *SubService) Delete(ctx context.Context, id uuid.UUID) error {

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("subscriptionId", id.String()).
			Msg("Delete subscription: repository error")
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	s.logger.Info().
		Str("subscriptionId", id.String()).
		Msg("Subscription deleted successfully")
	return nil
}

func (s *SubService) GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	s.logger.Debug().
		Str("subscriptionId", id.String()).
		Msg("Getting subscription by ID")

	s.logger.Debug().
		Str("subscriptionId", id.String()).
		Msg("Fetching subscription from repository")

	subscription, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("subscriptionId", id.String()).
			Msg("Failed to fetch subscription from repository")
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	s.logger.Info().
		Str("subscriptionId", subscription.Id.String()).
		Str("serviceName", subscription.ServiceName).
		Str("userId", subscription.UserId.String()).
		Int64("price", subscription.Price).
		Msg("Subscription retrieved successfully")
	return subscription, nil
}

func (s *SubService) GetAll(ctx context.Context, page, pageSize int64) (dto.GetAllResponse, error) {
	s.logger.Debug().
		Int64("page", page).
		Int64("pageSize", pageSize).
		Msg("Getting all subscriptions")

	subscriptions, totalCount, totalPages, err := s.repo.GetAll(ctx, page, pageSize)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("page", page).
			Int64("pageSize", pageSize).
			Msg("Failed to fetch subscriptions from repository")
		return dto.GetAllResponse{}, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	s.logger.Info().
		Int64("page", page).
		Int64("pageSize", pageSize).
		Int64("totalCount", totalCount).
		Int64("totalPages", totalPages).
		Int("subscriptionsCount", len(subscriptions)).
		Msg("Subscriptions retrieved successfully")

	response := dto.GetAllResponse{
		Data: subscriptions,
		Pagination: &dto.PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

func (s *SubService) CalculateTotalCost(ctx context.Context, req dto.CostCalculationQueryRequest) (int64, error) {
	s.logger.Debug().
		Interface("filters", req).
		Msg("Calculating total cost")

	filter := &filters.SubFilter{
		UserID:      &req.UserID,
		ServiceName: &req.ServiceName,
		From:        &req.From,
		To:          &req.To,
	}

	total, err := s.repo.SumSubscriptionsCost(ctx, filter)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to calculate total cost")
		return 0, err
	}

	s.logger.Info().
		Int64("total", total).
		Msg("Total cost calculated")
	return total, nil
}

func (s *SubService) applyPartialUpdate(existing *models.Subscription, request dto.UpdateSubscriptionRequest) *models.Subscription {
	updated := &models.Subscription{
		Id:        existing.Id,
		UserId:    existing.UserId,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// Обновляем только переданные поля
	if request.ServiceName != nil {
		updated.ServiceName = *request.ServiceName
	} else {
		updated.ServiceName = existing.ServiceName
	}

	if request.Price != nil {
		updated.Price = *request.Price
	} else {
		updated.Price = existing.Price
	}

	if request.StartDate != nil {
		updated.StartDate = *request.StartDate
	} else {
		updated.StartDate = existing.StartDate
	}

	if request.EndDate != nil {
		updated.EndDate = request.EndDate
	} else {
		updated.EndDate = existing.EndDate
	}

	return updated
}
