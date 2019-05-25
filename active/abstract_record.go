package active

import (
	"database/sql"
	"github.com/infarmasistemas/go-abstract-record/active/abstract"
	"github.com/infarmasistemas/go-abstract-record/active/check"
	"github.com/infarmasistemas/go-abstract-record/active/instance"
	"github.com/infarmasistemas/go-abstract-record/active/query/post_funcs"
)

type AbstractRecord struct {
	abstractOps abstract.AbstractOps
	checkOps    check.CheckOps
	instanceOps instance.InstanceOps
}

func NewAbstractRecord(object interface{}, objectArray interface{}, db *sql.DB, extraOptions ...interface{}) *AbstractRecord {
	abstract := AbstractRecord{}
	abstract.abstractOps.Prepare(object, objectArray, db, extraOptions...)
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
	return e.checkOps.TreatEntry(e.instanceOps.NewSingle, e.instanceOps.NewArray, values...)
}

func (e *AbstractRecord) Save() error {
	return e.abstractOps.Save()
}

func (e *AbstractRecord) Delete() error {
	return e.abstractOps.Delete()
}

func (e *AbstractRecord) Update(values ...interface{}) error {
	return e.checkOps.TreatValuesForUpdate(e.abstractOps.Update, values...)
}

func (e *AbstractRecord) Paginate(pk string, offset, pageSize int) *AbstractRecord {
	if pageSize > 0 {
		query, values := post_funcs.PaginationFunc(pk, offset, pageSize)
		e.abstractOps.SqlOps.GetComposition().GetComposer().PostQuery.AddPostQuery(query)
		e.abstractOps.SqlOps.GetComposition().GetComposer().PostQuery.AddValues(values...)
	}

	return e
}
