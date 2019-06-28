package collection_of_attributes

import (
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"github.com/kataras/iris/core/errors"
	"reflect"
)

type Attribute struct {
	Table       string
	ValueOf     reflect.Value
	StructField reflect.StructField
	Options     options.OptionsOps
}

func NewAttribute(table string, valueOf reflect.Value, field reflect.StructField, extraOptions ...interface{}) Attribute {
	return Attribute{table, valueOf, field, options.NewOptionsOps(extraOptions...)}
}

func (a *Attribute) garTagWithoutLimits() string {
	if gar, presence := a.StructField.Tag.Lookup("gar"); presence {
		return gar
	}
	return ""
}

func (a *Attribute) garTagWithLimits() string {
	if gar, presence := a.StructField.Tag.Lookup("gar"); presence {
		if a.Options.FieldsPresentInQueryCustomFields(a.Table, gar) {
			return gar
		}
	}
	return ""
}

func (a *Attribute) GarTag() string {
	if a.Options.CheckIfCustomFieldsAreFromThisTable(a.Table) && a.Options.QueryCustomFieldsPresent() {
		return a.garTagWithLimits()
	} else {
		return a.garTagWithoutLimits()
	}
}

func (a *Attribute) JSONTag() string {
	if jsonTag, presence := a.StructField.Tag.Lookup("json"); presence {
		return jsonTag
	}
	return ""
}

func (a *Attribute) Interface() (interface{}, error) {
	if !a.ValueOf.Elem().IsValid() {
		return nil, errors.New("can't call interface on this value")
	}

	return a.ValueOf.Elem().Interface(), nil
}

func (a *Attribute) RelTagPresent() bool {
	_, presence := a.StructField.Tag.Lookup("rel")

	return presence
}

func (a *Attribute) AddrInterface() interface{} {
	//if a.Options.CheckIfCustomFieldsAreFromThisTable(a.Table) && a.Options.QueryCustomFieldsPresent() {
	//	if a.Options.FieldsPresentInQueryCustomFields(a.Table, a.StructField.Name) {
	//		return a.ValueOf.Addr().Interface()
	//	}
	//
	//	return nil
	//}
	//
	//if !a.RelTagPresent() {
	return a.ValueOf.Addr().Interface()
	//}

	//return nil
}

func (a *Attribute) DefaultValue() (string, bool) {
	return a.StructField.Tag.Lookup("default_value")
}
