package composer

import "strings"

type Where struct {
	base       string
	conditions []string
}

func NewWhere() Where {
	where := Where{}
	where.base = "WHERE "

	return where
}

func (w *Where) Valid() bool {
	if len(w.conditions) > 0 {
		return true
	} else {
		return false
	}
}

func (w *Where) AddCondition(value ...string) {
	w.conditions = append(w.conditions, value...)
}

func (w *Where) Build() string {
	var sb strings.Builder

	if !w.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(w.base)
	for index, condition := range w.conditions {
		if index == (len(w.conditions) - 1) {
			sb.WriteString(condition)
			sb.WriteString(" = ?")
		} else {
			sb.WriteString(condition)
			sb.WriteString(" = ?")
			sb.WriteString(" AND ")
		}
	}

	sb.WriteString(" ")

	return sb.String()
}
