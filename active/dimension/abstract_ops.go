package dimension

import (
	"database/sql"
	"github.com/infarmasistemas/goorm/active/dimension/query"
)

type AbstractOps struct {
	SqlOps      query.SQLOps
}

func (d *AbstractOps) prepare(object interface{}, objectArray interface{}, db *sql.DB) {
	d.SqlOps = *query.NewSQLOps(object, objectArray, db)
}

func (d *AbstractOps) Find(values ...interface{}) error {
	return d.SqlOps.Select(values...)
}

func (d *AbstractOps) Where(values ...interface{}) error {
	return d.SqlOps.Where(values...)
}