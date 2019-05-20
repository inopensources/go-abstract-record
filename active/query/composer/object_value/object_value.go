package object_value

import (
	"reflect"
)

type ObjectValue struct {
	object interface{}
}

func NewObjectValue(object interface{}) ObjectValue {
	return ObjectValue{object}
}

func (o *ObjectValue) IsValid() bool {
	if o.object == nil {
		return false
	}

	if !o.isPointer() {
		return true
	}

	if reflect.ValueOf(o.object).Elem().IsValid() {
		return true
	}

	return false
}

func (o *ObjectValue) ReturnSQL() string {
	if o.IsValid() {
		return " = ?"
	}

	return " IS NULL"
}

func (o *ObjectValue) GetObject() interface{} {
	if o.isPointer() {
		return reflect.ValueOf(o.object).Elem().Interface()
	}

	return o.object
}

func (o *ObjectValue) isPointer() bool {
	if reflect.ValueOf(o.object).Kind() == reflect.Ptr {
		return true
	}

	return false
}
