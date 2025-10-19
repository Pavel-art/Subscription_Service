package dto

import (
	"SubscriptionService/internal/core/models"
)

type GetAllResponse struct {
	Data       []*models.Subscription `json:"data"`
	Pagination *PaginationInfo        `json:"pagination"`
}

type PaginationInfo struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
}
