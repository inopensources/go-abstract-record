package composer

import "strings"

type Count struct {
	base    string
	columns []string
}

func NewCount() Count {
	count := Count{}
	count.base = "SELECT COUNT "

	return count
}

func (c *Count) AddColumn(value ...string) {
	c.columns = append(c.columns, value...)
}

func (c *Count) Valid() bool {
	if len(c.columns) > 0 {
		return true
	} else {
		return false
	}
}

func (c *Count) Build() string {
	var sb strings.Builder

	if !c.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(c.base)
	sb.WriteString("(")
	for index, column := range c.columns {
		if index == (len(c.columns)-1) {
			sb.WriteString(column)
		} else {
			sb.WriteString(column)
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	sb.WriteString(" ")

	return sb.String()
}
