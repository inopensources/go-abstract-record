package query

import "strings"

type QueriesOps struct {
	preQueries  []string
	midQueries  []string
	postQueries []string
	values      []interface{}
}

func (q *QueriesOps) AddPreQuery(preQuery string) {
	q.preQueries = append(q.preQueries, preQuery)
}

func (q *QueriesOps) AddMidQuery(midQuery string) {
	q.midQueries = append(q.midQueries, midQuery)
}

func (q *QueriesOps) AddPostQuery(postQuery string) {
	q.postQueries = append(q.postQueries, postQuery)
}

func (q *QueriesOps) returnBuiltQueryAndValues() (string, []interface{}) {
	var sb strings.Builder
	//Building first part of the query
	sb.WriteString(q.buildQueryFromPreQueries())

	//Building second part of the query
	sb.WriteString(q.buildQueryFromMidQueries())

	//Building thirdpart of the query
	sb.WriteString(q.buildQueryFromPostQueries())

	sb.WriteString(";")

	return sb.String(), q.values
}

func (q *QueriesOps) buildQueryFromPreQueries() string {
	return q.buildQueryFromParts(q.preQueries)
}

func (q *QueriesOps) buildQueryFromMidQueries() string {
	return q.buildQueryFromParts(q.midQueries)
}

func (q *QueriesOps) buildQueryFromPostQueries() string {
	return q.buildQueryFromParts(q.postQueries)
}

func (q *QueriesOps) buildQueryFromParts(values []string) string {
	var sb strings.Builder
	for _, partialQuery := range values {
		sb.WriteString(partialQuery)
	}

	return sb.String()
}

func (q *QueriesOps) AddValues(values ...interface{}) {
	q.values = append(q.values, values...)
}
