package dto

type PaginationQuery struct {
	Page   int
	Limit  int
	Offset int
}

type PaginatedList[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
