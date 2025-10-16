package utils

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	PreviousPage int `json:"previous_page"`
	NextPage     int `json:"next_page"`
}

func BuildPagination(page, pageSize, totalRecords int) Pagination {
	if page <= 0 || pageSize <= 0 || totalRecords <= 0 {
		return Pagination{}
	}

	totalPages := totalRecords / pageSize
	if totalRecords%pageSize != 0 {
		totalPages++
	}

	previousPage := page - 1
	if previousPage < 1 {
		previousPage = 1
	} else if previousPage > totalPages {
		previousPage = totalPages
	}

	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	return Pagination{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		PreviousPage: previousPage,
		NextPage:     nextPage,
	}
}
