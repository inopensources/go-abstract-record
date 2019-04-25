package dimension

import (
	"database/sql"
)

type AbstractRecord struct {
	checkLayer  Check
	abstractOps AbstractOps
}

func NewAbstractRecord(object interface{}, objectArray interface{}, db *sql.DB) *AbstractRecord {
	abstract := AbstractRecord{}
	abstract.abstractOps.prepare(object, objectArray, db)
	return &abstract
}

func (e *AbstractRecord) Find(values ...interface{}) error {
	return e.checkLayer.CheckAndExecute(e.abstractOps.Find, values...)
}

func (e *AbstractRecord) Where(values ...interface{}) error {
	return e.checkLayer.CheckAndExecute(e.abstractOps.Where, values...)
}
