package composer

import (
	"github.com/infarmasistemas/go-abstract-record/active/query/collection_of_attributes"
	"strings"
)

type Join struct {
	base        string
	parentTable string
	prefix      string
	childTables []collection_of_attributes.TablePkFk
}

func NewJoin() Join {
	join := Join{}
	join.base = "LEFT OUTER JOIN "

	return join
}

func (s *Join) SetParentTable(parentTable string) {
	s.parentTable = parentTable
}

func (s *Join) SetPrefix(somePrefix string) {
	s.prefix = somePrefix
}
func (s *Join) SetChildTables(childTables ...collection_of_attributes.TablePkFk) {
	s.childTables = append(s.childTables, childTables...)
}

func (s *Join) Valid() bool {
	if len(s.childTables) > 0 {
		return true
	} else {
		return false
	}
}

func (s *Join) Build() string {
	var sb strings.Builder

	if !s.Valid() {
		return ""
	}

	// Writing tables
	for _, childTable := range s.childTables {
		sb.WriteString(s.base)
		sb.WriteString(childTable.TableChild)
		sb.WriteString(" ON ")

		for indexPk, Pk := range childTable.Pk {
			if indexPk == (len(childTable.Pk) - 1) {
				sb.WriteString(s.prefix + childTable.TableChild + "." + Pk)
				sb.WriteString(" = ")
				sb.WriteString(s.prefix + s.parentTable + "." + childTable.Fk[indexPk])
			} else {
				sb.WriteString(s.prefix + childTable.TableChild + "." + Pk)
				sb.WriteString(" = ")
				sb.WriteString(s.prefix + s.parentTable + "." + childTable.Fk[indexPk])
				sb.WriteString(" AND ")
			}
		}
		sb.WriteString(" ")
	}
	sb.WriteString(" ")

	return sb.String()
}
