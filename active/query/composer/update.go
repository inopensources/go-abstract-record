package composer

import "strings"

type Update struct {
	base    string
	tableName string
}

func NewUpdate() Update {
	update := Update{}
	update.base = "UPDATE "

	return update
}

func (u *Update) Valid() bool {
	if len(u.tableName) > 0 {
		return true
	} else {
		return false
	}
}

func (u *Update) AddTableName(value string) {
	u.tableName = value
}

func (u *Update) Build() string {
	var sb strings.Builder

	if !u.Valid() {
		return ""
	}
	sb.WriteString(u.base)
	sb.WriteString(u.tableName)
	sb.WriteString(" ")

	return sb.String()
}
