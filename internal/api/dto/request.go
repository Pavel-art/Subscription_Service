package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required,min=2,max=100"`
	Price       int64      `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID  `json:"user_id" binding:"required,uuid"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *int64     `json:"price,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type CostCalculationQueryRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
}
