package has_one

import (
	"reflect"
)

type RelatedField struct {
	structField reflect.StructField
	field       reflect.Value
	fk          interface{}
}

func NewRelatedField(structField reflect.StructField, field reflect.Value, fk interface{}) RelatedField {
	return RelatedField{
		structField: structField,
		field:       field,
		fk:          fk,
	}
}

func (r *RelatedField) getParams() []reflect.Value {
	var params []reflect.Value

	params = append(params, reflect.ValueOf(r.structField.Tag.Get("through")))
	params = append(params, reflect.ValueOf(r.fk))

	return params
}

func (r *RelatedField) FetchRelation(object interface{}) {
	//Reflect object's value
	someObject := reflect.ValueOf(object).Elem()

	//Capture object field
	targetField := someObject.FieldByName(r.structField.Tag.Get("has_one"))

	//Capture type of object
	someType := reflect.TypeOf(targetField.Interface()).Elem()

	//Create new object of said type
	maNew := reflect.New(someType)

	//Capture method called Find (it should exist)
	method := reflect.ValueOf(maNew.Interface()).MethodByName("Find")

	//Call method find passing parameters
	method.Call(r.getParams())

	//Write object to field specified via the through tag
	targetField.Set(maNew)
}