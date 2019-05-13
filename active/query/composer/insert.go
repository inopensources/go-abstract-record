package composer

import "strings"

type Insert struct {
	base    string
	columns []string
	tableName string
}

func NewInsert() Insert {
	insert := Insert{}
	insert.base = "INSERT INTO "

	return insert
}

func (i *Insert) Valid() bool {
	if len(i.columns) > 0 {
		return true
	} else {
		return false
	}
}

func (i *Insert) AddColumn(value ...string) {
	i.columns = append(i.columns, value...)
}

func (i *Insert) AddTableName(value string) {
	i.tableName = value
}

func (i *Insert) Build() string {
	var sb strings.Builder

	if !i.Valid() {
		return ""
	}

	// Writing tables
	sb.WriteString(i.base)
	sb.WriteString(i.tableName)
	sb.WriteString(" ")
	sb.WriteString("(")
	for index, column := range i.columns {
		if index == (len(i.columns)-1) {
			sb.WriteString(column)
		} else {
			sb.WriteString(column)
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	sb.WriteString(" ")
	sb.WriteString("VALUES")
	sb.WriteString(" ")
	sb.WriteString("(")
	for index := range i.columns {
		if index == (len(i.columns)-1) {
			sb.WriteString("?")
		} else {
			sb.WriteString("?")
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	return sb.String()
}
