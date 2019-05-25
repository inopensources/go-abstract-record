package relationships

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldInterface struct {
	sliceOfThrough []string
	sliceOfVia     []string
	FieldName      string
	values         []interface{}
}

func NewFieldInterface(object interface{}, field reflect.StructField, relationshipType string) FieldInterface {
	var fieldInterface FieldInterface
	fieldInterface.FieldName = field.Tag.Get(relationshipType)
	fieldInterface.sliceOfThrough = strings.Split(field.Tag.Get("through"), ":")
	fieldInterface.sliceOfVia = strings.Split(field.Tag.Get("via"), ":")

	if len(fieldInterface.sliceOfThrough) > 1 {
		fmt.Println("Multiple FKs")
	} else {
		fmt.Println("Single FK")
	}

	fieldInterface.getValuesFromObject(object)

	return fieldInterface
}

func (f *FieldInterface) getValuesFromObject(object interface{}) {
	objectReflection := reflect.ValueOf(object).Elem()

	for _, via := range f.sliceOfVia {
		f.values = append(f.values, objectReflection.FieldByName(via).Elem().Interface())
	}
}

func (f *FieldInterface) GetParams() (params []reflect.Value) {
	for index, columnName := range f.sliceOfThrough {
		params = append(params, reflect.ValueOf(columnName))
		params = append(params, reflect.ValueOf(f.values[index]))
	}

	return params
}
