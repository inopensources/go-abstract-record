package collection_of_attributes

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"reflect"
	"strings"
)

type CollectionOfAttributes struct {
	Table                  string
	CollectionOfAttributes []Attribute
}

func NewCollectionOfAttributes(object interface{}, extraOptions ...interface{}) CollectionOfAttributes {
	table := strings.Split(reflect.TypeOf(object).String(), ".")[len(strings.Split(reflect.TypeOf(object).String(), "."))-1]

	s := reflect.ValueOf(object).Elem()
	typeOfT := s.Type()

	var attributes []Attribute
	for i := 0; i < s.NumField(); i++ {
		attributes = append(attributes, NewAttribute(table, s.Field(i), typeOfT.Field(i), extraOptions...))
	}

	return CollectionOfAttributes{table, attributes}
}

func (a *CollectionOfAttributes) PointersToAttributes() (pointerList []interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if addrInterface := attribute.AddrInterface(); addrInterface != nil {
			pointerList = append(pointerList, addrInterface)
		}
	}

	return pointerList
}

func (a *CollectionOfAttributes) AttributeValuesAsSlice() (pointerList []interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if !attribute.RelTagPresent() {
			pointerList = append(pointerList, attribute.ValueOf.Interface())
		}
	}

	return pointerList
}

func (a *CollectionOfAttributes) AttributeValuesAsArray() (sliceOfValues []interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if value, err := attribute.Interface(); err == nil {
			sliceOfValues = append(sliceOfValues, value)
		} else {
			sliceOfValues = append(sliceOfValues, nil)
		}
	}

	return sliceOfValues
}

func (a *CollectionOfAttributes) AttributesAsColumnNames() (columns []string) {
	for _, attribute := range a.CollectionOfAttributes {
		if garTag := attribute.GarTag(); garTag != "" {
			columns = append(columns, garTag)
		}
	}

	return columns
}

func (a *CollectionOfAttributes) RealColName(colName string) (string, error) {
	for _, attribute := range a.CollectionOfAttributes {
		if colName == attribute.JSONTag() {
			if attribute.RelTagPresent() {
				return "", errors.New("column is present, but is a relationship")
			}

			return attribute.GarTag(), nil
		}
	}

	return "", errors.New("column not present")
}

//Remove this in the future
func (a *CollectionOfAttributes) Conditions(values ...interface{}) (conditions []string) {
	var queryValues []interface{}

	for i, condition := range values {
		if i%2 == 0 {
			conditions = append(conditions, fmt.Sprintf("%s", condition))
		} else {
			queryValues = append(queryValues, condition)
		}
	}

	return conditions
}
