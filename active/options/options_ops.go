package options

import (
	"github.com/infarmasistemas/go-abstract-record/active/helpers"
	"reflect"
)

const DeepQuery = 0
const MaxLevel = 1
const CurrentLevel = 2
const QueryCustomFields = 3

type OptionsOps struct {
	// True: Every collection_of_attributes is rendered with all of its relationships
	// False: Objects are rendered without their relationships
	DeepQuery         bool
	MaxLevel          int
	CurrentLevel      int
	QueryCustomFields helpers.Limit
}

//Using variadic functions here, as it would be a bummer to pass
//far too many values.
//Index | Param
//  0   | DeepQuery
//  1   | MaxLevel
func NewOptionsOps(extraOptions ...interface{}) OptionsOps {
	if extraOptions != nil {
		return OptionOpsFromExtraOptions(extraOptions...)
	}

	return OptionsOps{
		DeepQuery:         false,
		CurrentLevel:      0,
		MaxLevel:          1,
		QueryCustomFields: helpers.NewLimit(),
	}
}

func OptionOpsFromExtraOptions(extraOptions ...interface{}) OptionsOps {
	var newOptionOps OptionsOps

	newOptionOps.setDeepQueryFromInterface(extraOptions)
	newOptionOps.setMaxLevelFromInterface(extraOptions)
	newOptionOps.setCurrentLevelFromInterface(extraOptions)
	newOptionOps.setQueryCustomFieldsFromInterface(extraOptions)

	return newOptionOps
}

func (o *OptionsOps) setDeepQueryFromInterface(extraOptions []interface{}) {
	if o.validDepth(DeepQuery, extraOptions) {
		o.DeepQuery = extraOptions[DeepQuery].(bool)
	}
}

func (o *OptionsOps) setMaxLevelFromInterface(extraOptions []interface{}) {
	if o.validDepth(MaxLevel, extraOptions) {
		o.MaxLevel = extraOptions[MaxLevel].(int)
	}
}

func (o *OptionsOps) setCurrentLevelFromInterface(extraOptions []interface{}) {
	if o.validDepth(CurrentLevel, extraOptions) {
		o.CurrentLevel = extraOptions[CurrentLevel].(int)
	}
}

func (o *OptionsOps) setQueryCustomFieldsFromInterface(extraOptions []interface{}) {
	if o.validDepth(QueryCustomFields, extraOptions) {
		o.QueryCustomFields = extraOptions[QueryCustomFields].(helpers.Limit)
	}
}

func (o *OptionsOps) validDepth(level int, extraOptions []interface{}) bool {
	if len(extraOptions) > level {
		return true
	}

	return false
}

func (o *OptionsOps) getOptionsAsArray() []interface{} {
	var values []interface{}
	s := reflect.ValueOf(o).Elem()

	for i := 0; i < s.NumField(); i++ {
		values = append(values, s.Field(i).Interface())
	}

	return values
}

func (o *OptionsOps) GetOptionsAsArrayOfValues() (optionsAsValues []reflect.Value) {
	for _, value := range o.getOptionsAsArray() {
		optionsAsValues = append(optionsAsValues, reflect.ValueOf(value))
	}

	return optionsAsValues
}

func (o *OptionsOps) IncreaseCurrentLevel() {
	o.CurrentLevel += 1
}

func (o *OptionsOps) CheckIfCurrentLevelBiggerThanMaxLevel() bool {
	if o.CurrentLevel > o.MaxLevel {
		return true
	}

	return false
}

func (o *OptionsOps) CheckIfCustomFieldsAreFromThisTable(tableName string) bool {
	if o.QueryCustomFields.GetLimits()[tableName] != nil {
		return true
	}

	return false
}

func (o *OptionsOps) QueryCustomFieldsPresent() bool {
	if o.QueryCustomFields.Valid() {
		return true
	}

	return false
}

func (o *OptionsOps) QueryCustomFieldsAsSlice(tableName string) []string {
	return o.QueryCustomFields.GetLimits()[tableName]
}

func (o *OptionsOps) FieldsPresentInQueryCustomFields(tableName, fieldName string) bool {
	for _, value := range o.QueryCustomFieldsAsSlice(tableName) {
		if value == fieldName {
			return true
		}
	}

	return false
}
