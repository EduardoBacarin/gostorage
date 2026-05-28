package helpers

import "math"

type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPages int64 `json:"total_pages"`
}

type PaginationConfig struct {
	Limit int64
	Skip  int64
}

func PreparePagination(page, limit int64) PaginationConfig {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	return PaginationConfig{
		Limit: limit,
		Skip:  (page - 1) * limit,
	}
}

func NewPaginatedResult[T any](data []T, total, page, limit int64) *PaginatedResult[T] {
	if limit < 1 {
		limit = 10
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	if data == nil {
		data = []T{}
	}

	return &PaginatedResult[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
