package composer

import (
	"fmt"
	"strings"
)

type Max struct {
	base   string
	column string
}

func NewMax() Max {
	count := Max{}
	count.base = "SELECT MAX "

	return count
}

func (c *Max) AddColumn(value interface{}) {
	c.column = fmt.Sprintf("%s", value)
}

func (c *Max) Valid() bool {
	if c.column != "" {
		return true
	} else {
		return false
	}
}

func (c *Max) Build() string {
	var sb strings.Builder

	if !c.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(c.base)
	sb.WriteString("(")

	sb.WriteString(c.column)

	sb.WriteString(")")
	sb.WriteString(" ")

	return sb.String()
}
