package instance

import (
	json2 "encoding/json"
	"github.com/infarmasistemas/go-abstract-record/active/query/collection_of_attributes"
)

type InstanceOps struct {
	CollectionOfAttributes collection_of_attributes.CollectionOfAttributes
	object                 interface{}
	objectArray            interface{}
}

func (i *InstanceOps) Prepare(object interface{}, objectArray interface{}) {
	i.object = object
	i.objectArray = objectArray
	i.CollectionOfAttributes = collection_of_attributes.NewCollectionOfAttributes(object)
}

// TODO: Code a way to instantiate objects using new
func (i *InstanceOps) NewSingle(json []byte) error {
	return json2.Unmarshal(json, &i.object)
}

func (i *InstanceOps) NewArray(json []byte) error {
	return json2.Unmarshal(json, &i.objectArray)
}

func (i *InstanceOps) NewWithDefaults() error {
	return i.CollectionOfAttributes.SetDefaultValues()
}
