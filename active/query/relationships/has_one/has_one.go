package has_one

import (
	"reflect"
)

type HasOne struct {
	relatedFields []RelatedField
	object        interface{}
}

func NewHasOne(object interface{}) HasOne {
	return HasOne{
		object: object,
	}
}

func (h *HasOne) AddRelatedField(structField reflect.StructField, field reflect.Value) {
	h.relatedFields = append(h.relatedFields, NewRelatedField(structField, field))
}

func (h *HasOne) RelatedFieldsPresent() bool {
	if len(h.relatedFields) > 0 {
		return true
	}

	return false
}

func (h *HasOne) CheckPresenceOfHasOneRelationship(field reflect.StructField) (string, bool) {
	return field.Tag.Lookup("has_one")
}

func (h *HasOne) FetchRelatedObjects() {
	for _, relatedField := range h.relatedFields {
		relatedField.FetchRelation(h.object)
	}
}
