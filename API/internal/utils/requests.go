package utils

import (
	"net/http"
	"strconv"
)

type UserInformation struct {
	UserId string
	Email  string
	TaxId  string
}

func ExtractPaginationParams(r *http.Request) (pageNum, pageSizeNum int64) {
	var err error
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	pageNum, err = strconv.ParseInt(page, 10, 0)
	if err != nil {
		pageNum = 0
	}
	pageSizeNum, err = strconv.ParseInt(pageSize, 10, 0)
	if err != nil {
		pageSizeNum = 0
	}
	return
}
