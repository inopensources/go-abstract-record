package relationships

import "reflect"

type ObjectInterface struct {
	object interface{}
}

func NewObjectInterface(object interface{}) ObjectInterface {
	return ObjectInterface{
		object: object,
	}
}

func (o *ObjectInterface) ObjectValue() reflect.Value {
	return reflect.ValueOf(o.object).Elem()
}

func (o *ObjectInterface) FieldByName(fieldName string) reflect.Value {
	return o.ObjectValue().FieldByName(fieldName)
}

func (o *ObjectInterface) TypeOfField(fieldName string) reflect.Type {
	return reflect.TypeOf(o.FieldByName(fieldName).Interface()).Elem()
}

func (o *ObjectInterface) NewObjectFromFieldNameType(fieldName string) reflect.Value {
	return reflect.New(o.TypeOfField(fieldName))
}

func (o *ObjectInterface) NewSliceFromFieldType(fieldName string) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(o.TypeOfField(fieldName)), 0, 1)
}

func (o *ObjectInterface) SetFieldValueByName(fieldName string, value reflect.Value) {
	o.FieldByName(fieldName).Set(value)
}

func (o *ObjectInterface) NewSliceFromFieldTypeForSlice(fieldName string) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(o.TypeOfField(fieldName).Elem()), 0, 1)
}

func (o *ObjectInterface) NewObjectFromFieldNameTypeForSlices(fieldName string) reflect.Value {
	return reflect.New(o.TypeOfField(fieldName).Elem())
}

func (o *ObjectInterface) NewSliceFromFieldTypeForSlices(fieldName string) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(o.TypeOfField(fieldName).Elem()), 0, 1)
}
