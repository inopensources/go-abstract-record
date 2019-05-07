package instance

import (
	json2 "encoding/json"
)

type InstanceOps struct {
	object      interface{}
	objectArray interface{}
}

func (i *InstanceOps) Prepare(object interface{}, objectArray interface{}) {
	i.object = object
	i.objectArray = objectArray
}

// TODO: Code a way to instantiate objects using new
func (i *InstanceOps) NewSingle(json []byte) error {
	return json2.Unmarshal(json, &i.object)
}

func (i *InstanceOps) NewArray(json []byte) error {
	return json2.Unmarshal(json, &i.objectArray)
}