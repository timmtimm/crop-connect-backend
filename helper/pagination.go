package helper

import (
	"errors"
	"marketplace-backend/util"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaginationParam struct {
	Page  string
	Limit string
	Sort  string
	Order string
}

type QueryPagination struct {
	Skip  int64
	Limit int64
	Sort  string
	Order int
}

func PaginationToQuery(c echo.Context, availableSort []string) (QueryPagination, error) {
	pagination := PaginationParam{
		Page:  c.Param("page"),
		Limit: c.QueryParam("limit"),
		Sort:  c.QueryParam("sort"),
		Order: c.QueryParam("order"),
	}

	if pagination.Limit == "" {
		pagination.Limit = "10"
	}

	if pagination.Sort == "" {
		pagination.Sort = "createdAt"
	}

	if pagination.Order == "" {
		pagination.Order = "desc"
	}

	page, err := strconv.Atoi(pagination.Page)
	if err != nil {
		return QueryPagination{}, errors.New("halaman harus berupa angka")
	} else if page < 1 {
		return QueryPagination{}, errors.New("halaman tidak boleh kurang dari 1")
	}

	limit, err := strconv.Atoi(pagination.Limit)
	if err != nil {
		return QueryPagination{}, errors.New("limit harus berupa angka")
	} else if limit < 1 {
		return QueryPagination{}, errors.New("limit tidak boleh kurang dari 1")
	}

	if !util.CheckStringOnArray(availableSort, pagination.Sort) {
		return QueryPagination{}, errors.New("penyortiran tidak tersedia")
	}

	var convertOrder int
	if !util.CheckStringOnArray([]string{"asc", "desc"}, pagination.Order) {
		return QueryPagination{}, errors.New("urutan hanya bisa asc dan desc")
	} else if pagination.Order == "asc" {
		convertOrder = 1
	} else {
		convertOrder = -1
	}

	return QueryPagination{
		Skip:  int64((page - 1) * limit),
		Limit: int64(limit),
		Sort:  pagination.Sort,
		Order: convertOrder,
	}, nil
}

func ConvertToPaginationResponse(query QueryPagination, totalData int) Page {
	return Page{
		Size:        int(query.Limit),
		TotalData:   totalData,
		CurrentPage: int(query.Skip/query.Limit) + 1,
		TotalPage:   int(totalData/int(query.Limit)) + 1,
	}
}
