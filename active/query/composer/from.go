package composer

import "strings"

type From struct {
	base 	string
	tables  []string
}

func NewFrom() From {
	from := From{}
	from.base = "FROM "

	return from
}

func (f *From) AddTableName(value ...string) {
	f.tables = append(f.tables, value...)
}

func (f *From) Valid() bool {
	if len(f.tables) > 0 {
		return true
	} else {
		return false
	}
}

func (f *From) Build() string {
	var sb strings.Builder

	if !f.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(f.base)
	for index, table := range f.tables {
		if index == (len(f.tables)-1) {
			sb.WriteString(table)
		} else {
			sb.WriteString(table)
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" ")

	return sb.String()
}
