package query

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type SQLOps struct {
	composition CompositionOps
	Object      interface{}
	ObjectArray interface{}
	Db          *sql.DB
}

func NewSQLOps(object interface{}, objectArray interface{}, db *sql.DB) *SQLOps {
	sqlOps := SQLOps{}
	sqlOps.composition = *NewCompositionOps(object)
	sqlOps.Object = object
	sqlOps.ObjectArray = objectArray
	sqlOps.Db = db

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

	fmt.Println(result.LastInsertId())

	return nil
}

func (s *SQLOps) GetComposition() *CompositionOps {
	return &s.composition
}