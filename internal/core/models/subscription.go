package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrServiceNameRequired = errors.New("service name is required")
	ErrPriceInvalid        = errors.New("price must be positive")
	ErrStartDateRequired   = errors.New("start date is required")
	ErrEndDateBeforeStart  = errors.New("end date cannot be before start date")
)

type Subscription struct {
	Id          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int64      `json:"price"`
	UserId      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (s *Subscription) Validate() error {
	if s.ServiceName == "" {
		return ErrServiceNameRequired
	}

	if s.Price <= 0 {
		return ErrPriceInvalid
	}

	if s.StartDate.IsZero() {
		return ErrStartDateRequired
	}

	// Валидация опционального EndDate
	if s.EndDate != nil && !s.EndDate.IsZero() {
		if s.EndDate.Before(s.StartDate) {
			return ErrEndDateBeforeStart
		}
	}

	return nil
}

func NewSubscription(
	serviceName string,
	price int64,
	userID uuid.UUID,
	startDate time.Time,
	endDate *time.Time) (*Subscription, error) {
	sub := &Subscription{
		Id:          uuid.New(),
		ServiceName: serviceName,
		Price:       price,
		UserId:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := sub.Validate(); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *Subscription) IsActive() bool {
	now := time.Now()

	if s.EndDate == nil || s.EndDate.IsZero() {
		return true
	}

	return now.After(s.StartDate) && now.Before(*s.EndDate)
}
