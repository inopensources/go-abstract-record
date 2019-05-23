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

func NewRelationshipOps(object interface{}) RelationshipOps {
	return RelationshipOps{
		object: object,
		hasOne: has_one.NewHasOne(object),
		hasMany: has_many.NewHasMany(object),
	}
}

func (r *RelationshipOps) checkForRelationships() {
	s := reflect.ValueOf(r.object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		field := typeOfT.Field(i)
		if value, presenceHasOne := r.hasOne.CheckPresenceOfHasOneRelationship(field); presenceHasOne {
			fmt.Println("Has one", value)
			r.hasOne.AddRelatedField(field, s.Field(i))
		}

		if value, presenceHasMany := r.hasMany.CheckPresenceOfHasManyRelationship(field); presenceHasMany {
			fmt.Println("Has many", value)
			r.hasMany.AddRelatedField(field, s.Field(i))
		}
	}
}

func (r *RelationshipOps) fetchHasOneRelatedObjects() {
	if r.hasOne.RelatedFieldsPresent() {
		r.hasOne.FetchRelatedObjects()
	}
}

func (r *RelationshipOps) fetchHasManyRelatedObjects() {
	if r.hasMany.RelatedFieldsPresent() {
		r.hasMany.FetchRelatedObjects()
	}
}
