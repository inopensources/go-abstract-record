package collection_of_attributes

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"github.com/kataras/iris/core/errors"
	"reflect"
	"strconv"
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

func (a *Attribute) SetValueWithReflection(someValue interface{}) {
	reflectedValue := reflect.ValueOf(someValue)
	pt := reflect.PtrTo(reflectedValue.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())               // create a reflect.Value of type *T.
	pv.Elem().Set(reflectedValue)              // sets pv to point to underlying value of v.

	a.ValueOf.Set(pv)
}

func (a *Attribute) TryToParseDefaultValueAndSetDefault() error {
	switch a.ValueOf.Type().String() {
	case "*string":
		return a.tryToParseDefaultValueAsStringAndSetValueWithReflection()
	case "*bool":
		return a.tryToParseDefaultValueAsBoolAndSetValueWithReflection()
	case "*float32":
		return a.tryToParseDefaultValueAsFloat32AndSetValueWithReflection()
	case "*float64":
		return a.tryToParseDefaultValueAsFloat64AndSetValueWithReflection()
	case "*int":
		return a.tryToParseDefaultValueAsIntAndSetValueWithReflection()
	case "*int32":
		return a.tryToParseDefaultValueAsInt32AndSetValueWithReflection()
	case "*int64":
		return a.tryToParseDefaultValueAsInt64AndSetValueWithReflection()
	default:
		return errors.New(fmt.Sprintf("Cannot deal with type %s. It can be added by request.", a.ValueOf.Type().String()))
	}
}

func (a *Attribute) tryToParseDefaultValueAsStringAndSetValueWithReflection() error {
	defaultValue, _ := a.DefaultValue()
	a.SetValueWithReflection(defaultValue)

	return nil
}

func (a *Attribute) tryToParseDefaultValueAsBoolAndSetValueWithReflection() error {
	defaultValue, _ := a.DefaultValue()

	val, err := strconv.ParseBool(defaultValue)
	if err != nil {
		return err
	}

	a.SetValueWithReflection(val)

	return nil
}

func (a *Attribute) tryToParseDefaultValueAsFloat32AndSetValueWithReflection() error {
	defaultValue, _ := a.DefaultValue()

	val, err := strconv.ParseFloat(defaultValue, 4)
	if err != nil {
		return err
	}

	a.SetValueWithReflection(float32(val))

	return nil
}

func (a *Attribute) tryToParseDefaultValueAsFloat64AndSetValueWithReflection() error {
	defaultValue, _ := a.DefaultValue()

	val, err := strconv.ParseFloat(defaultValue, 4)
	if err != nil {
		return err
	}

	a.SetValueWithReflection(val)

	return nil
}

func (a *Attribute) tryToParseDefaultValueAsIntAndSetValueWithReflection() error {
	defaultValue, _ := a.DefaultValue()

	val, err := strconv.ParseInt(defaultValue, 10, 4)
	if err != nil {
		return err
	}

	a.SetValueWithReflection(int(val))

	return nil
}

func (a *Attribute) tryToParseDefaultValueAsInt32AndSetValueWithReflection() error {
	defaultValue, _ := a.DefaultValue()

	val, err := strconv.ParseInt(defaultValue, 10, 4)
	if err != nil {
		return err
	}

	a.SetValueWithReflection(int32(val))

	return nil
}

func (a *Attribute) tryToParseDefaultValueAsInt64AndSetValueWithReflection() error {
	defaultValue, _ := a.DefaultValue()

	val, err := strconv.ParseInt(defaultValue, 10, 4)
	if err != nil {
		return err
	}

	a.SetValueWithReflection(int64(val))

	return nil
}
