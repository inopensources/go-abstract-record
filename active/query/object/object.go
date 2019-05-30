package object

import (
	"errors"
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"reflect"
	"strings"
)

type Object struct {
	object           interface{}
	Table            string
	AttributesGar    []string
	AttributesJson   map[string]string
	AttributesValues []interface{}
	QueryValues      []interface{}
	PointerList      []interface{}
	options          options.OptionsOps
}

func NewObject(object interface{}, extraOptions ...interface{}) Object {
	var newObject Object
	newObject.object = object
	newObject.options = options.NewOptionsOps(extraOptions...)
	newObject.discoverTable()
	newObject.discoverAttributesAndpointerList()

	return newObject
}

func (o *Object) discoverTable() {
	o.Table = strings.Split(reflect.TypeOf(o.object).String(), ".")[len(strings.Split(reflect.TypeOf(o.object).String(), "."))-1]
}

func (o *Object) discoverAttributesAndpointerList() {
	if o.options.CheckIfCustomFieldsAreFromThisTable(o.Table) && o.options.QueryCustomFieldsPresent() {
		o.buildPointerListAndAttributesWithLimits()
	} else {
		o.buildPointerListAndAttributes()
	}
}

func (o *Object) parseJsonGar(field reflect.StructField) (string, string) {
	var finalGar, finalJson string

	if _, presence := field.Tag.Lookup("rel"); presence {
		return finalGar, finalJson
	}

	finalJson = field.Tag.Get("json")
	finalGar = field.Tag.Get("gar")

	return finalGar, finalJson
}

func (o *Object) AttributesAsColumnNames() []string {
	var columns []string
	for _, attributeName := range o.AttributesGar {
		columns = append(columns, attributeName)
	}

	return columns
}

func (o *Object) AttributeValuesAsArray() []interface{} {
	var values []interface{}
	for _, value := range o.AttributesValues {
		if reflect.ValueOf(value).Elem().IsValid() {
			values = append(values, reflect.ValueOf(value).Elem().Interface())
		} else {
			values = append(values, nil)
		}
	}

	return values
}

func (o *Object) Conditions(values ...interface{}) []string {
	var queryValues []interface{}

	var conditions []string

	for i, condition := range values {
		if i%2 == 0 {
			conditions = append(conditions, fmt.Sprintf("%s", condition))
		} else {
			queryValues = append(queryValues, condition)
		}
	}

	o.QueryValues = queryValues
	return conditions
}

func (o *Object) buildPointerListAndAttributesWithLimits() {
	var fieldList []interface{}
	var attributeValues []interface{}
	var attributeListGar []string

	o.AttributesJson = make(map[string]string, 0)

	s := reflect.ValueOf(o.object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		finalGar, finalJson := o.parseJsonGar(typeOfT.Field(i))
		if len(finalGar) > 0 && len(finalJson) > 0 {
			someField := s.Field(i)
			if o.options.FieldsPresentInQueryCustomFields(o.Table, finalGar) {
				fieldList = append(fieldList, someField.Addr().Interface())
				attributeListGar = append(attributeListGar, finalGar)
				o.AttributesJson[finalJson] = finalGar

				attributeValues = append(attributeValues, someField.Interface())
			}
		}
	}

	o.AttributesGar = attributeListGar
	o.PointerList = fieldList
	o.AttributesValues = attributeValues
}

func (o *Object) buildPointerListAndAttributes() {
	var fieldList []interface{}
	var attributeValues []interface{}
	var attributeListGar []string

	o.AttributesJson = make(map[string]string, 0)

	s := reflect.ValueOf(o.object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		finalGar, finalJson := o.parseJsonGar(typeOfT.Field(i))

		if len(finalGar) > 0 && len(finalJson) > 0 {
			someField := s.Field(i)

			fieldList = append(fieldList, someField.Addr().Interface())
			attributeListGar = append(attributeListGar, finalGar)
			o.AttributesJson[finalJson] = finalGar

			attributeValues = append(attributeValues, someField.Interface())
		}
	}

	o.AttributesGar = attributeListGar
	o.PointerList = fieldList
	o.AttributesValues = attributeValues
}

func (o *Object) GetRealColName(value string) (string, error) {
	if len(o.AttributesJson[value]) == 0 {
		return "", errors.New("field not found")
	}

	return o.AttributesJson[value], nil
}
