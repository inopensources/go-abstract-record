package composer

import "strings"

type Delete struct {
	base    string
	invoked bool
}

func NewDelete() Delete {
	delete := Delete{}
	delete.base = "DELETE "

	return delete
}

func (d *Delete) Call() {
	d.invoked = true
}

func (d *Delete) Valid() bool {
	return d.invoked
}

func (d *Delete) Build() string {
	var sb strings.Builder

	if !d.Valid() {
		return ""
	}
	sb.WriteString(d.base)

	return sb.String()
}
