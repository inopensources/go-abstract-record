package tail_funcs

import "fmt"

func PaginationFunc(pk string, offset, pageSize int) func() string {
	return NewPagination(pk, offset, pageSize).Paginate
}

type Pagination struct {
	Pk       string
	Offset   int
	PageSize int
}

func NewPagination(pk string, offset, pageSize int) *Pagination {
	pagination := Pagination{}
	pagination.Pk = pk
	pagination.Offset = offset
	pagination.PageSize = pageSize
	return &pagination
}

func (p *Pagination) Paginate() string {
	return fmt.Sprintf("ORDER BY %s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", p.Pk, p.Offset, p.PageSize)
}
