package dimension

import (
	"database/sql"
	"github.com/infarmasistemas/go-abstract-record/active/dimension/abstract"
	"github.com/infarmasistemas/go-abstract-record/active/dimension/check"
	"github.com/infarmasistemas/go-abstract-record/active/dimension/instance"
)

type AbstractRecord struct {
	checkOps    check.CheckOps
	instanceOps instance.InstanceOps
	abstractOps abstract.AbstractOps
}

func NewAbstractRecord(object interface{}, objectArray interface{}, db *sql.DB,  extraFuncs ...func() string) *AbstractRecord {
	abstract := AbstractRecord{}
	abstract.abstractOps.Prepare(object, objectArray, db, extraFuncs...)
	abstract.instanceOps.Prepare(object, objectArray)
	return &abstract
}

func (e *AbstractRecord) All() error {
	return e.abstractOps.All()
}

func (e *AbstractRecord) Find(values ...interface{}) error {
	return e.checkOps.CheckAndExecute(e.abstractOps.Find, values...)
}

func (e *AbstractRecord) Where(values ...interface{}) error {
	return e.checkOps.CheckAndExecute(e.abstractOps.Where, values...)
}

func (e *AbstractRecord) New(values ...interface{}) error {
	return  e.checkOps.TreatEntry(e.instanceOps.NewSingle, e.instanceOps.NewArray, values...)
}

func (e *AbstractRecord) Save() error {
	return e.abstractOps.Save()
}
