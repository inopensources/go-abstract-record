package composer

import "strings"

type Set struct {
	base    string
	columns []string
}

func NewSet() Set {
	set := Set{}
	set.base = "SET "

	return set
}

func (s *Set) AddColumn(value ...string) {
	s.columns = append(s.columns, value...)
}

func (s *Set) Valid() bool {
	if len(s.columns) > 0 {
		return true
	} else {
		return false
	}
}

func (s *Set) Build() string {
	var sb strings.Builder

	if !s.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(s.base)
	for index, column := range s.columns {
		if index == (len(s.columns)-1) {
			sb.WriteString(column)
			sb.WriteString(" = ?")
		} else {
			sb.WriteString(column)
			sb.WriteString(" = ?")
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" ")

	return sb.String()
}
