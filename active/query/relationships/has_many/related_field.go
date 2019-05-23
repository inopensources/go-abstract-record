package has_many

import (
	"reflect"
)

type RelatedField struct {
	structField reflect.StructField
	field       reflect.Value
}

func NewRelatedField(structField reflect.StructField, field reflect.Value) RelatedField {
	return RelatedField{
		structField: structField,
		field:       field,
	}
}

func (r *RelatedField) getParams(object interface{}) []reflect.Value {
	var params []reflect.Value

	//Get with value
	fk := reflect.ValueOf(object).Elem().FieldByName(r.structField.Tag.Get("with")).Interface()

	params = append(params, reflect.ValueOf(0))
	params = append(params, reflect.ValueOf(0))
	params = append(params, reflect.ValueOf(r.structField.Tag.Get("through")))
	params = append(params, reflect.ValueOf(fk))

	return params
}

func (r *RelatedField) FetchRelations(object interface{}) {
	//Reflect object's value
	someObject := reflect.ValueOf(object).Elem()

	//Capture object field
	targetField := someObject.FieldByName(r.structField.Tag.Get("has_many"))

	//Capture type of object
	someType := reflect.TypeOf(targetField.Interface()).Elem()

	//Create new object of said type
	maNew := reflect.New(someType)

	//Capture method called Find (it should exist)
	method := reflect.ValueOf(maNew.Interface()).MethodByName("All")

	//Call method find passing parameters
	returnedValues := method.Call(r.getParams(object))

	//Write returned values to field specified via the through tag
	targetField.Set(returnedValues[0])
}
