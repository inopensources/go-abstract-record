package dimension

import (
	"database/sql"
	"github.com/infarmasistemas/goorm/active/dimension/abstract"
	"github.com/infarmasistemas/goorm/active/dimension/check"
	"github.com/infarmasistemas/goorm/active/dimension/instance"
)

type AbstractRecord struct {
	checkOps    check.CheckOps
	instanceOps instance.InstanceOps
	abstractOps abstract.AbstractOps
}

func NewAbstractRecord(object interface{}, objectArray interface{}, db *sql.DB) *AbstractRecord {
	abstract := AbstractRecord{}
	abstract.abstractOps.Prepare(object, objectArray, db)
	abstract.instanceOps.Prepare(object, objectArray)
	return &abstract
}

func (e *AbstractRecord) Find(values ...interface{}) error {
	return e.checkOps.CheckAndExecute(e.abstractOps.Find, values...)
}

func (e *AbstractRecord) Where(values ...interface{}) error {
	return e.checkOps.CheckAndExecute(e.abstractOps.Where, values...)
}

func (e *AbstractRecord) New(values ...interface{}) error {
	return  e.checkOps.CheckAndExecute(e.instanceOps.New, values...)
}

func (e *AbstractRecord) Save() error {
	return e.abstractOps.Save()
}
