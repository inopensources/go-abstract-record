package query

import (
	"database/sql"
	"encoding/json"
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
	OptionsOps    options.OptionsOps
}

func NewSQLOps(object interface{}, objectArray interface{}, db *sql.DB, extraOptions ...interface{}) *SQLOps {
	sqlOps := SQLOps{}
	sqlOps.composition = *NewCompositionOps(object, extraOptions...)
	sqlOps.Object = object
	sqlOps.ObjectArray = objectArray
	sqlOps.Db = db
	sqlOps.OptionsOps = options.NewOptionsOps(extraOptions...)

	return &sqlOps
}

func (s *SQLOps) Select(values ...interface{}) error {
	query, queryValues := s.composition.Select(values...)

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryValues...)
	if err != nil {
		return err
	}

	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		colNames, err := rows.Columns()

		err = rows.Scan(s.composition.CollectionOfAttributes.PointersToAttributesFromColumnNames(colNames...)...)
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
	//going to be eager loaded until MaxLevel is reached
	if s.OptionsOps.DeepQuery {
		//Method below checks if we're already too deep into the relationship /:
		if !s.OptionsOps.CheckIfCurrentLevelBiggerThanMaxLevel() {
			s.dealWithRelationships()
		}
	}

	return err
}

func (s *SQLOps) Where(values ...interface{}) error {
	query, queryValues := s.composition.Select(values...)

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryValues...)
	if err != nil {
		return err
	}

	defer rows.Close()

	resultCount := 0

	for rows.Next() {
		colNames, err := rows.Columns()

		err = rows.Scan(s.composition.CollectionOfAttributes.PointersToAttributesFromColumnNames(colNames...)...)
		if err != nil {
			fmt.Println(err)
		}

		//If DeepQuery is set to true, the relationships are
		//going to be loaded
		//ATTENTION: This may slowdown your query, as a new SQL
		//query will be created for every parent collection_of_attributes returned
		if s.OptionsOps.DeepQuery {
			//Method below checks if we're already too deep into the relationship /:
			if !s.OptionsOps.CheckIfCurrentLevelBiggerThanMaxLevel() {
				s.dealWithRelationships()
			}
		}

		// Improve this
		objType := reflect.TypeOf(s.Object).Elem()
		newObj := reflect.New(objType)
		jsonOld, err := json.Marshal(s.Object)
		json.Unmarshal(jsonOld, newObj.Interface())

		valuePtr := reflect.ValueOf(s.ObjectArray)
		value := valuePtr.Elem()
		value.Set(reflect.Append(value, newObj.Elem()))

		resultCount++
	}

	return err
}

func (s *SQLOps) Insert() error {
	query, queryValues := s.composition.Insert()

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(queryValues...)
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

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(queryValues...)
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

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(queryValues...)
	if err != nil {
		return err
	}

	recordsAffected, _ := result.RowsAffected()
	fmt.Println(recordsAffected, "records affected")

	return nil
}

func (s *SQLOps) Count(values ...interface{}) (int, error) {
	var count int
	query, queryValues := s.composition.Count(values...)

	fmt.Println("QUERY:", query)
	fmt.Println("VALUES:", queryValues)

	rows := s.Db.QueryRow(query)

	err := rows.Scan(&count)

	return count, err
}

func (s *SQLOps) GetComposition() *CompositionOps {
	return &s.composition
}

func (s *SQLOps) dealWithRelationships() {
	//return
	deepOptions := s.OptionsOps
	deepOptions.IncreaseCurrentLevel()

	//In this case, we're checking if the level of the options
	//we're passing down the rabbit hole is still valid
	if deepOptions.CheckIfCurrentLevelBiggerThanMaxLevel() {
		return
	}

	//Creating new relationship collection_of_attributes
	s.relationships = NewRelationshipOps(s.Object, deepOptions)

	//Method below checks if the current collection_of_attributes has got relationships
	s.relationships.checkForRelationships()

	//Deep query deactivated
	//if s.relationships.hasOne.RelatedFieldsPresent() {
	//	s.relationships.fetchHasOneRelatedObjects()
	//}

	if s.relationships.hasMany.RelatedFieldsPresent() {
		s.relationships.fetchHasManyRelatedObjects()
	}
}
