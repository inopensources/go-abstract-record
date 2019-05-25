package relationships

import (
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"reflect"
)

type RelationshipInterface struct {
	relatedFields []RelatedField
	object        interface{}
	options       options.OptionsOps
}

func NewRelationshipInterface(object interface{}, options options.OptionsOps) RelationshipInterface {
	return RelationshipInterface{
		object:  object,
		options: options,
	}
}

func (r *RelationshipInterface) AddRelatedField(fieldInterface FieldInterface) {
	r.relatedFields = append(r.relatedFields, NewRelatedField(fieldInterface))
}

func (r *RelationshipInterface) RelatedFieldsPresent() bool {
	if len(r.relatedFields) > 0 {
		return true
	}

	return false
}

func (r *RelationshipInterface) CheckPresenceOfHasOneRelationship(field reflect.StructField) (string, bool) {
	return field.Tag.Lookup("has_one")
}

func (r *RelationshipInterface) CheckPresenceOfHasManyRelationship(field reflect.StructField) (string, bool) {
	return field.Tag.Lookup("has_many")
}

func (r *RelationshipInterface) FetchRelatedObjectsForHasOne() {
	if r.options.CheckIfCurrentLevelBiggerThanMaxLevel() {
		return
	}

	for _, relatedField := range r.relatedFields {
		relatedField.FetchRelationForSingleValue(r.object, r.options)
	}
}

func (r *RelationshipInterface) FetchRelatedObjectsForHasMany() {
	if r.options.CheckIfCurrentLevelBiggerThanMaxLevel() {
		return
	}

	for _, relatedField := range r.relatedFields {
		relatedField.FetchRelationForSliceValue(r.object, r.options)
	}
}
