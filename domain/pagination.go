package domain

type PaginationMetadata struct {
	TotalCount   int64 `json:"total_count"`
	TotalPages   int64 `json:"total_pages"`
	CurrentPage  int64 `json:"current_page"`
	ItemsPerPage int64 `json:"items_per_page"`
}
