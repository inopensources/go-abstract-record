package composer

import (
	"github.com/infarmasistemas/go-abstract-record/active/query/composer/object_value"
	"strings"
)

type Where struct {
	base         string
	conditions   []string
	objectValues []object_value.ObjectValue
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

func (w *Where) AddValues(values ...interface{}) {
	for _, value := range values {
		w.objectValues = append(w.objectValues, object_value.NewObjectValue(value))
	}
}

func (w *Where) Build() string {
	var sb strings.Builder

	if !w.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(w.base)
	for index, condition := range w.conditions {
		sb.WriteString(condition)

		if w.indexSecondToLast(index) {
			if w.objectValuesPresent() {
				if w.containsLike(condition) {
					sb.WriteString(" ?")
					break
				}

				sb.WriteString(w.objectValues[index].ReturnSQL())
			} else {
				sb.WriteString(" = ?")
			}
			break
		}

		if w.objectValuesPresent() {
			if w.containsLike(condition) {
				sb.WriteString(" ?")
			} else {
				sb.WriteString(w.objectValues[index].ReturnSQL())
			}
		} else {
			sb.WriteString(" = ?")
		}

		sb.WriteString(" AND ")
	}

	sb.WriteString(" ")

	return sb.String()
}

func (w *Where) objectValuesPresent() bool {
	if len(w.objectValues) > 0 {
		return true
	}

	return false
}

func (w *Where) getValues() []interface{} {
	var values []interface{}
	for _, value := range w.objectValues {
		if value.IsValid() {
			values = append(values, value.GetObject())
		}
	}

	return values
}

func (w *Where) containsLike(condition string) bool {
	if strings.HasSuffix(condition, " LIKE") {
		return true
	}

	return false
}

func (w *Where) indexSecondToLast(index int) bool {
	if index == (len(w.conditions) - 1) {
		return true
	}

	return false
}