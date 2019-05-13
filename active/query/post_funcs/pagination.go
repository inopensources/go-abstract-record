package post_funcs

import "fmt"

func PaginationFunc(pk string, offset, pageSize int) (string, []interface{})  {
	return NewPagination(pk, offset, pageSize).Paginate()
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

func (p *Pagination) Paginate() (string, []interface{}) {
	var values []interface{}

	values = append(values, p.Offset)
	values = append(values, p.PageSize)

	//Since the Order By expects the name of a column instead of a value, it has to be
	//replaced directly into the string instead as part of a prepared statement.
	query := fmt.Sprintf("ORDER BY %s OFFSET ? ROWS FETCH NEXT ? ROWS ONLY", p.Pk)

	return query, values
}
