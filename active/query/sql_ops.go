package query

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"reflect"
)

type SQLOps struct {
	relationships RelationshipOps
	composition   CompositionOps
	Object        interface{}
	ObjectArray   interface{}
	Db            *sql.DB
	optionsOps    options.OptionsOps
}

func NewSQLOps(object interface{}, objectArray interface{}, db *sql.DB, extraOptions ...bool) *SQLOps {
	sqlOps := SQLOps{}
	sqlOps.composition = *NewCompositionOps(object)
	sqlOps.Object = object
	sqlOps.ObjectArray = objectArray
	sqlOps.Db = db
	sqlOps.relationships = *NewRelationshipOps(object)
	sqlOps.optionsOps = options.NewOptionsOps(extraOptions...)

	return &sqlOps
}

func (s *SQLOps) Select(values ...interface{}) error {
	query, queryValues := s.composition.Select(values...)

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	rows, err := s.Db.Query(query, queryValues...)
	if err != nil {
		return err
	}

	resultCount := 0
	for rows.Next() {
		err := rows.Scan(s.composition.pointerList...)
		if err != nil {
			return err
		}

		if resultCount > 1 {
			return errors.New("multiple results")
		}

		resultCount++
	}

	if resultCount == 0 {
		return errors.New("No record found")
	}

	//If DeepQuery is set to true, the relationships are
	//going to be loaded
	if s.optionsOps.DeepQuery {
		//Method below checks if the current object has got relationships
		s.relationships.checkForRelationships(s.Object)
	}

	return err
}

func (s *SQLOps) Where(values ...interface{}) error {
	query, queryValues := s.composition.Select(values...)

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	rows, err := s.Db.Query(query, queryValues...)
	if err != nil {
		return err
	}

	resultCount := 0
	for rows.Next() {
		err := rows.Scan(s.composition.pointerList...)
		if err != nil {
			fmt.Println(err)
		}

		valuePtr := reflect.ValueOf(s.ObjectArray)
		value := valuePtr.Elem()
		value.Set(reflect.Append(value, reflect.ValueOf(s.Object).Elem()))

		resultCount++
	}

	return err
}

func (s *SQLOps) Insert() error {
	query, queryValues := s.composition.Insert()

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	result, err := s.Db.Exec(query, queryValues...)
	if err != nil {
		return err
	}

	recordsAffected, _ := result.RowsAffected()
	fmt.Println(recordsAffected, "records affected")

	return nil
}

func (s *SQLOps) Delete() error {
	query, queryValues := s.composition.Delete()

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	result, err := s.Db.Exec(query, queryValues...)
	if err != nil {
		return err
	}

	recordsAffected, _ := result.RowsAffected()
	fmt.Println(recordsAffected, "records affected")

	return err
}

func (s *SQLOps) Update(values ...interface{}) error {
	query, queryValues := s.composition.Update(values...)

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	result, err := s.Db.Exec(query, queryValues...)
	if err != nil {
		return err
	}

	recordsAffected, _ := result.RowsAffected()
	fmt.Println(recordsAffected, "records affected")

	return nil
}

func (s *SQLOps) GetComposition() *CompositionOps {
	return &s.composition
}
