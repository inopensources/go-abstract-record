package collection_of_attributes

import (
	"reflect"
	"strings"
)

type TablePkFk struct {
	TableParent string
	TableChild  string
	Pk          []string
	Fk          []string
}

func NewTablePkFk(tableParent, tableChild string, field reflect.StructField) TablePkFk {
	var tablePkFk TablePkFk

	tablePkFk.TableParent = tableParent
	tablePkFk.TableChild = tableChild

	tablePkFk.Pk = strings.Split(field.Tag.Get("through"), ":")
	tablePkFk.Fk = strings.Split(field.Tag.Get("via"), ":")

	return tablePkFk
}
