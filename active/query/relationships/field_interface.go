package relationships

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldInterface struct {
	sliceOfThrough []string
	sliceOfVia     []string
	mapOfOnly      map[string][]string
	FieldName      string
	values         []interface{}
}

func NewFieldInterface(object interface{}, field reflect.StructField, relationshipType string) FieldInterface {
	var fieldInterface FieldInterface

	fieldInterface.FieldName = field.Tag.Get(relationshipType)
	fieldInterface.sliceOfThrough = strings.Split(field.Tag.Get("through"), ":")
	fieldInterface.sliceOfVia = strings.Split(field.Tag.Get("via"), ":")

	//if _, customFieldsPresent := field.Tag.Lookup("only"); customFieldsPresent {
	//	objectInterface := NewObjectInterface(collection_of_attributes)
	//	relatedObject := objectInterface.NewObjectFromFieldNameType(fieldInterface.FieldName).Interface()
	//
	//	someMap := make(map[string][]string, 0)
	//	someMap[object2.NewObject(relatedObject).Table] = strings.Split(field.Tag.Get("only"), ":")
	//	fieldInterface.mapOfOnly = someMap
	//}

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
		if objectReflection.FieldByName(via).Elem().IsValid() {
			f.values = append(f.values, objectReflection.FieldByName(via).Elem().Interface())
		}
	}
}

func (f *FieldInterface) GetParams() (params []reflect.Value, valid bool) {
	for index, columnName := range f.sliceOfThrough {
		if len(f.values) > index {
			params = append(params, reflect.ValueOf(columnName))
			params = append(params, reflect.ValueOf(f.values[index]))
		}
	}

	if len(params) > 0 {
		return params, true
	}

	return params, false
}
