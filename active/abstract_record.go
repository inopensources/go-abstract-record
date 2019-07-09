package active

import (
	"database/sql"
	"github.com/infarmasistemas/go-abstract-record/active/abstract"
	"github.com/infarmasistemas/go-abstract-record/active/check"
	"github.com/infarmasistemas/go-abstract-record/active/instance"
	"github.com/infarmasistemas/go-abstract-record/active/query/post_funcs"
	"github.com/infarmasistemas/go-abstract-record/active/utils/logger"
)

type AbstractRecord struct {
	abstractOps abstract.AbstractOps
	checkOps    check.CheckOps
	instanceOps instance.InstanceOps
}

func NewAbstractRecord(object interface{}, objectArray interface{}, db *sql.DB, extraOptions ...interface{}) *AbstractRecord {
	logger.NewLogger().Info("STARTING ENGINES")
	abstract := AbstractRecord{}
	abstract.abstractOps.Prepare(object, objectArray, db, extraOptions...)
	abstract.instanceOps.Prepare(object, objectArray)
	return &abstract
}

func (e *AbstractRecord) All() error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.abstractOps.All()
}

func (e *AbstractRecord) Count(values ...interface{}) (int, error) {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.abstractOps.Count(values...)
}

func (e *AbstractRecord) Max(value interface{}) (int, error) {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.abstractOps.Max(value)
}

func (e *AbstractRecord) Find(values ...interface{}) error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.checkOps.CheckAndExecute(e.abstractOps.Find, values...)
}

func (e *AbstractRecord) Where(values ...interface{}) error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.checkOps.CheckAndExecute(e.abstractOps.Where, values...)
}

func (e *AbstractRecord) New(values ...interface{}) error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.checkOps.TreatEntry(e.instanceOps.NewSingle, e.instanceOps.NewArray, values...)
}

func (e *AbstractRecord) NewWithDefaults(values ...interface{}) error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.instanceOps.NewWithDefaults()
}

func (e *AbstractRecord) Save() error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.abstractOps.Save()
}

func (e *AbstractRecord) Delete() error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

	return e.abstractOps.Delete()
}

func (e *AbstractRecord) Update(values ...interface{}) error {
	defer logger.NewLogger().Info("SHUTTING OFF ENGINES")

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

func (e *AbstractRecord) Inner() *AbstractRecord {
	options := e.abstractOps.SqlOps.OptionsOps
	object := e.abstractOps.SqlOps.Object
	objectArray := e.abstractOps.SqlOps.ObjectArray
	db := e.abstractOps.SqlOps.Db

	options.Inner = true

	e.abstractOps.Prepare(object, objectArray, db, options.GetOptionsAsSliceOfInterface()...)

	return e
}

func (e *AbstractRecord) SQL(sqlQuery string, values ...interface{}) ([]map[string]interface{}, error) {
	return e.abstractOps.SQL(sqlQuery, values...)
}
