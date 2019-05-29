package relationships

import (
	"github.com/infarmasistemas/go-abstract-record/active/helpers"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"reflect"
)

type RelatedField struct {
	fieldInterface FieldInterface
	options        options.OptionsOps
}

func NewRelatedField(fieldInterface FieldInterface) RelatedField {
	return RelatedField{
		fieldInterface: fieldInterface,
	}
}

func (r *RelatedField) getPrepareParams(objectInterface ObjectInterface, options options.OptionsOps) []reflect.Value {
	var params []reflect.Value
	params = append(params, objectInterface.NewSliceFromFieldType(r.fieldInterface.FieldName))
	params = append(params, options.GetOptionsAsArrayOfValues()...)

	return params
}

func (r *RelatedField) FetchRelationForSingleValue(object interface{}, options options.OptionsOps) {
	objectInterface := NewObjectInterface(object)

	values, valid := r.fieldInterface.GetParams()
	if !valid {
		return
	}

	//Create new object of said type
	maNew := objectInterface.NewObjectFromFieldNameType(r.fieldInterface.FieldName)

	if len(r.fieldInterface.mapOfOnly) > 0 {
		options.QueryCustomFields = helpers.NewLimitFromMap(r.fieldInterface.mapOfOnly)
	}

	params := r.getPrepareParams(objectInterface, options)

	//Call method find passing parameters
	prepareValue := maNew.MethodByName("Prepare").Call(params)

	findMethod := prepareValue[0].MethodByName("Find")

	findMethod.Call(values)

	//Write object to field specified via the through tag
	objectInterface.SetFieldValueByName(r.fieldInterface.FieldName, maNew)
}

func (r *RelatedField) FetchRelationForSliceValue(object interface{}, options options.OptionsOps) {
	objectInterface := NewObjectInterface(object)

	//The guard clause below checks if a given parent
	//object contains valid values to be used as params
	// for the child query
	values, valid := r.fieldInterface.GetParams()
	if !valid {
		return
	}

	fieldName := r.fieldInterface.FieldName

	//Create new object of said type
	maNew := objectInterface.NewObjectFromFieldNameTypeForSlices(fieldName)

	originalSlice := objectInterface.NewSliceFromFieldTypeForSlices(fieldName)

	sliceOfType := reflect.New(originalSlice.Type())
	sliceOfType.Elem().Set(originalSlice)

	var params []reflect.Value
	params = append(params, sliceOfType)
	params = append(params, options.GetOptionsAsArrayOfValues()...)

	//Call method find passing parameters
	prepareValue := maNew.MethodByName("Prepare").Call(params)

	findMethod := prepareValue[0].MethodByName("Where")

	findMethod.Call(values)

	//Write object to field specified via the through tag
	objectInterface.SetFieldValueByName(fieldName, sliceOfType)
}
