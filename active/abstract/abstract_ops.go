package abstract

import (
	"database/sql"
	"github.com/infarmasistemas/go-abstract-record/active/query"
)

type AbstractOps struct {
	SqlOps query.SQLOps
}

func (d *AbstractOps) Prepare(object interface{}, objectArray interface{}, db *sql.DB, extraFuncs ...func() string) {
	d.SqlOps = *query.NewSQLOps(object, objectArray, db, extraFuncs...)
}

func (d *AbstractOps) All() error {
	return d.SqlOps.Where()
}

func (d *AbstractOps) Find(values ...interface{}) error {
	return d.SqlOps.Select(values...)
}

func (d *AbstractOps) Where(values ...interface{}) error {
	return d.SqlOps.Where(values...)
}

func (d *AbstractOps) Save(values ...interface{}) error {
	return d.SqlOps.Insert()
}
