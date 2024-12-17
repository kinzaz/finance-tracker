package dto

type PaginationRequestDto struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PaginationResponseDto[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"total_count"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
}
