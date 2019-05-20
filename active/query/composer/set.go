package composer

import (
	"github.com/infarmasistemas/go-abstract-record/active/query/composer/object_value"
	"strings"
)

type Set struct {
	base         string
	columns      []string
	objectValues []object_value.ObjectValue
}

func NewSet() Set {
	set := Set{}
	set.base = "SET "

	return set
}

func (s *Set) AddColumn(value ...string) {
	s.columns = append(s.columns, value...)
}

func (s *Set) AddValues(values ...interface{}) {
	for _, value := range values {
		s.objectValues = append(s.objectValues, object_value.NewObjectValue(value))
	}
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
		if index == (len(s.columns) - 1) {
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

func (s *Set) objectValuesPresent() bool {
	if len(s.objectValues) > 0 {
		return true
	}

	return false
}

func (s *Set) getValues() []interface{} {
	var values []interface{}
	for _, value := range s.objectValues {
		values = append(values, value.GetObject())

	}

	return values
}
