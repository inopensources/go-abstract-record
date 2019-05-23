package has_many

import "reflect"

type HasMany struct {
	relatedFields []RelatedField
	object        interface{}
}

func NewHasMany(object interface{}) *HasMany {
	return &HasMany{
		object: object,
	}
}

func (h *HasMany) AddRelatedField(structField reflect.StructField, field reflect.Value) {
	h.relatedFields = append(h.relatedFields, NewRelatedField(structField, field))
}

func (h *HasMany) RelatedFieldsPresent() bool {
	if len(h.relatedFields) > 0 {
		return true
	}

	return false
}

func (h *HasMany) FetchRelatedObjects(object interface{}) {
	for _, relatedField := range h.relatedFields {
		relatedField.FetchRelations(object)
	}
}
