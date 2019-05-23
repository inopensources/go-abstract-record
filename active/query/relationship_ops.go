package query

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/query/relationships/has_many"
	"github.com/infarmasistemas/go-abstract-record/active/query/relationships/has_one"
	"reflect"
)

type RelationshipOps struct {
	object  interface{}
	hasOne  has_one.HasOne
	hasMany has_many.HasMany
}

func NewRelationshipOps(object interface{}) *RelationshipOps {
	return &RelationshipOps{
		object: object,
	}
}

func (r *RelationshipOps) checkForRelationships(object interface{}) {
	s := reflect.ValueOf(r.object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		field := typeOfT.Field(i)
		if value, presenceHasOne := field.Tag.Lookup("has_one"); presenceHasOne {
			fmt.Println("Has one", value)
			r.hasOne.AddRelatedField(field, s.Field(i), s.Field(i).Elem().Interface())
		}

		if value, presenceHasMany := field.Tag.Lookup("has_many"); presenceHasMany {
			fmt.Println("Has many", value)
			r.hasMany.AddRelatedField(field, s.Field(i))
		}
	}

	//If there are relationships for this model in particular
	if r.hasOne.RelatedFieldsPresent() {
		r.hasOne.FetchRelatedObjects(object)
	}

	if r.hasMany.RelatedFieldsPresent() {
		r.hasMany.FetchRelatedObjects(object)
	}
}

//func (r *RelationshipOps) dealWithHasOne(structField reflect.StructField, field reflect.Value) {
//
//	hasOne := structField.Tag.Get("has_one")
//	through := structField.Tag.Get("through")
//
//	fmt.Println("Field", field.Elem())
//	fmt.Println("Struct", structField)
//
//	//HasOne
//	fmt.Println("HasOne", hasOne)
//	fmt.Println("Thorugh", through)
//	fmt.Println("Fk", field.Elem().Interface())
//
//	//Treat object
//	someObject := reflect.ValueOf(r.object).Elem()
//
//	targetField := someObject.FieldByName(hasOne)
//
//	someType := reflect.TypeOf(targetField.Interface()).Elem()
//
//	maNew := reflect.New(someType)
//
//	fmt.Println(maNew.Type())
//
//	method := reflect.ValueOf(maNew.Interface()).MethodByName("Find")
//
//	fmt.Println(method)
//
//	in := make([]reflect.Value, method.Type().NumIn())
//
//	param := reflect.ValueOf(`{ "Codigo" : 308 }`)
//
//	for index := 0; index < method.Type().NumIn(); index++ {
//		in[index] = param
//	}
//
//	method.Call(in)
//
//	k := maNew.Elem().Interface()
//
//	fmt.Println(k)
//
//	//Write relationship to object struct
//
//	targetField.Set(maNew)
//
//	fmt.Println("Done")
//
//}
//
//
