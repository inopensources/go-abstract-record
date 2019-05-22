package has_one

import "reflect"

type HasOne struct {
	relatedFields []RelatedField
	object        interface{}
}

func NewHasOne(object interface{}) *HasOne {
	return &HasOne{
		object: object,
	}
}

func (h *HasOne) AddRelatedField(structField reflect.StructField, field reflect.Value, fk interface{}) {
	h.relatedFields = append(h.relatedFields, NewRelatedField(structField, field, fk))
}

func (h *HasOne) RelatedFieldsPresent() bool {
	if len(h.relatedFields) > 0 {
		return true
	}

	return false
}

func (h *HasOne) FetchRelatedObjects(object interface{}) {
	for _, relatedField := range h.relatedFields {
		relatedField.FetchRelation(object)
	}
}