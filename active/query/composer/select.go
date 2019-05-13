package composer

import "strings"

type Select struct {
	base    string
	columns []string
}

func NewSelect() Select {
	selec := Select{}
	selec.base = "SELECT "

	return selec
}

func (s *Select) AddColumn(value ...string) {
	s.columns = append(s.columns, value...)
}

func (s *Select) Valid() bool {
	if len(s.columns) > 0 {
		return true
	} else {
		return false
	}
}

func (s *Select) Build() string {
	var sb strings.Builder

	if !s.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(s.base)
	for index, column := range s.columns {
		if index == (len(s.columns)-1) {
			sb.WriteString(column)
		} else {
			sb.WriteString(column)
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" ")

	return sb.String()
}
