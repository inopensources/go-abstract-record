package query

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type SQLOps struct {
	composition Composition
	object      interface{}
	objectArray interface{}
	db          *sql.DB
}

func NewSQLOps(object interface{}, objectArray interface{}, db *sql.DB) *SQLOps {
	sqlOps := SQLOps{}
	sqlOps.composition = *NewComposition(object)
	sqlOps.object = object
	sqlOps.objectArray = objectArray
	sqlOps.db = db

	return &sqlOps
}

func (s *SQLOps) Select(values ...interface{}) error {
	query, queryValues := s.composition.Select(values...)

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	rows, err := s.db.Query(query, queryValues...)
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

	rows, err := s.db.Query(query, queryValues...)
	if err != nil {
		return err
	}

	resultCount := 0
	for rows.Next() {
		err := rows.Scan(s.composition.pointerList...)
		if err != nil {
			fmt.Println(err)
		}

		valuePtr := reflect.ValueOf(s.objectArray)
		value := valuePtr.Elem()
		value.Set(reflect.Append(value, reflect.ValueOf(s.object).Elem()))

		resultCount++
	}

	return err
}