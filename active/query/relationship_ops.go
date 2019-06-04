package query

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"github.com/infarmasistemas/go-abstract-record/active/query/relationships"
	"reflect"
)

type RelationshipOps struct {
	object  interface{}
	hasOne  relationships.RelationshipInterface
	hasMany relationships.RelationshipInterface
}

func NewRelationshipOps(object interface{}, options options.OptionsOps) RelationshipOps {
	return RelationshipOps{
		object:  object,
		hasOne:  relationships.NewRelationshipInterface(object, options),
		hasMany: relationships.NewRelationshipInterface(object, options),
	}
}

func (r *RelationshipOps) checkForRelationships() {
	s := reflect.ValueOf(r.object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		field := typeOfT.Field(i)

		//if value, presenceHasOne := r.hasOne.CheckPresenceOfHasOneRelationship(field); presenceHasOne {
		//	fmt.Println("Has one", value)
		//	r.hasOne.AddRelatedField(relationships.NewFieldInterface(r.object, field, "has_one"))
		//}

		if value, presenceHasMany := r.hasMany.CheckPresenceOfHasManyRelationship(field); presenceHasMany {
			fmt.Println("Has many", value)
			r.hasMany.AddRelatedField(relationships.NewFieldInterface(r.object, field, "has_many"))
		}
	}
}

func (r *RelationshipOps) fetchHasOneRelatedObjects() {
	r.hasOne.FetchRelatedObjectsForHasOne()
}

func (r *RelationshipOps) fetchHasManyRelatedObjects() {
	r.hasMany.FetchRelatedObjectsForHasMany()
}
