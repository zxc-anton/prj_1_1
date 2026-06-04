package types

import "strconv"

type PaginateParams struct {
	Limit  int `json:"limit" validate:"required,min=10,max=30"`
	Offset int `json:"offset" validate:"required,min=0"`
}

func ParsePaginateParams(limitStr, offsetStr string) (limit int, offset int) {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 25 {
		limit = 10
	}
	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	return
}
