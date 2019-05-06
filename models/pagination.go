package models

import (
	"fmt"
	"github.com/kataras/iris"
	"strconv"
	"strings"
)

type Pagination struct {
	Page        int64
	MorePerPage int64
	Where       string
	GroupBy     string
	OrderBy     string
}

func MountPagination(ctx iris.Context) Pagination {
	var p Pagination
	querys := strings.Split(ctx.Request().URL.RawQuery, "&")

	for _, l := range querys {
		query := strings.Split(l, "=")
		switch query[0] {
		case "page":
			p.Page, _ = strconv.ParseInt(query[1], 10, 64)
		case "more_per_page":
			p.MorePerPage, _ = strconv.ParseInt(query[1], 10, 64)
		case "where":
			p.Where = query[1]
		case "order_by":
			p.GroupBy = query[1]
		case "group_by":
			p.OrderBy = query[1]
		default:
			fmt.Println("Errado")
		}
	}

	return p
}

type All struct {
	Data []interface{}
	Count int
}