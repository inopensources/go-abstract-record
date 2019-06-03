package collection_of_attributes

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"github.com/infarmasistemas/go-abstract-record/active/query/relationships"
	"github.com/kataras/iris/core/errors"
	"reflect"
	"strings"
)

type CollectionOfAttributes struct {
	Table                  string
	TableChild             []TablePkFk
	CollectionOfAttributes []Attribute
	options                options.OptionsOps
}

func NewCollectionOfAttributes(object interface{}, extraOptions ...interface{}) CollectionOfAttributes {
	var CollectionOfAttributes CollectionOfAttributes

	value := reflect.ValueOf(object).Elem()
	CollectionOfAttributes.Table = value.Type().Name()
	CollectionOfAttributes.options = options.NewOptionsOps(extraOptions...)

	CollectionOfAttributes.InspectAndCollectAttributes(object, extraOptions...)

	return CollectionOfAttributes
}

func (a *CollectionOfAttributes) InspectAndCollectAttributes(object interface{}, extraOptions ...interface{}) {
	objectInterface := relationships.NewObjectInterface(object)
	//fmt.Print(objectInterface)

	value := reflect.ValueOf(object).Elem()
	typeOfT := value.Type()
	table := value.Type().Name()

	var attributes []Attribute

	for i := 0; i < value.NumField(); i++ {
		attributes = append(attributes, NewAttribute(table, value.Field(i), typeOfT.Field(i), extraOptions...))

		if a.options.DeepQuery {
			if _, present := typeOfT.Field(i).Tag.Lookup("rel"); present {
				//Instantiate a new object from a nil related pointer
				maNew := objectInterface.NewObjectFromFieldNameType(typeOfT.Field(i).Name)

				tableChild := maNew.Elem().Type().Name()

				a.TableChild = append(a.TableChild, NewTablePkFk(a.Table, tableChild, typeOfT.Field(i)))

				//Writing that instantiated object to the object being handled
				objectInterface.SetFieldValueByName(typeOfT.Field(i).Name, maNew)

				//Pass those values down the rabbit hole to inspect and collect
				a.InspectAndCollectAttributes(maNew.Interface(), extraOptions...)
			}
		}
	}

	a.CollectionOfAttributes = append(a.CollectionOfAttributes, attributes...)
}

func (a *CollectionOfAttributes) PointersToAttributes() (pointerList []interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if addrInterface := attribute.AddrInterface(); addrInterface != nil {
			pointerList = append(pointerList, addrInterface)
		}
	}

	return pointerList
}

func (a *CollectionOfAttributes) PointersToAttributesFromColumnNames(columnNames ...string) (pointerList []interface{}) {
	for _, columnName := range columnNames {
		//v := a.FilterByColumnName(columnName)
		//fmt.Print(v)
		//pointerList = append(pointerList, a.AttributeFromColumnName(columnName))
		pointerList = append(pointerList, a.FilterByColumnName(columnName))
	}

	return pointerList
}

func (a *CollectionOfAttributes) AttributeFromColumnName(columnName string) (pointerToAttribute interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if !attribute.RelTagPresent() && attribute.garTagWithoutLimits() == columnName {
			if addrInterface := attribute.AddrInterface(); addrInterface != nil {
				pointerToAttribute = addrInterface
			}
		}
	}

	return pointerToAttribute
}

func (a *CollectionOfAttributes) AttributeValueFromColumnName(columnName string) (attributeValue interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if !attribute.RelTagPresent() && attribute.garTagWithoutLimits() == columnName {
			if value, err := attribute.Interface(); err == nil {
				attributeValue = value
			}
		}
	}

	return attributeValue
}

func (a *CollectionOfAttributes) AttributeValuesAsSlice() (pointerList []interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if !attribute.RelTagPresent() {
			pointerList = append(pointerList, attribute.ValueOf.Interface())
		}
	}

	return pointerList
}

func (a *CollectionOfAttributes) AttributeValuesAsArray() (sliceOfValues []interface{}) {
	for _, attribute := range a.CollectionOfAttributes {
		if value, err := attribute.Interface(); err == nil {
			sliceOfValues = append(sliceOfValues, value)
		} else {
			sliceOfValues = append(sliceOfValues, nil)
		}
	}

	return sliceOfValues
}

func (a *CollectionOfAttributes) AttributesAsColumnNamesForSelect() (columns []string) {
	for _, attribute := range a.CollectionOfAttributes {
		if garTag := attribute.GarTag(); garTag != "" {
			columns = append(columns, attribute.Table+"."+garTag+" as ["+attribute.Table+"."+garTag+"]")
		}
	}

	return columns
}

func (a *CollectionOfAttributes) AttributesAsColumnNamesForDelete() (columns []string) {
	return a.AttributesAsColumnNamesForUpdate()
}

func (a *CollectionOfAttributes) AttributesAsColumnNamesForUpdate() (columns []string) {
	for _, attribute := range a.CollectionOfAttributes {
		if garTag := attribute.GarTag(); garTag != "" {

			//If cast is present...
			if _, presence := attribute.StructField.Tag.Lookup("cast"); presence {
				columns = append(columns, fmt.Sprintf("CAST(%s AS VARCHAR(MAX))", attribute.GarTag())) //INSERT
			} else {
				columns = append(columns, garTag) //INSERT
			}
		}
	}

	return columns
}

func (a *CollectionOfAttributes) AttributesAsColumnNamesForInsert() (columns []string) {
	for _, attribute := range a.CollectionOfAttributes {
		if garTag := attribute.GarTag(); garTag != "" {
			columns = append(columns, garTag) //INSERT
		}
	}

	return columns
}

func (a *CollectionOfAttributes) RealColName(colName string) (string, error) {
	for _, attribute := range a.CollectionOfAttributes {
		if colName == attribute.JSONTag() {
			if attribute.RelTagPresent() {
				return "", errors.New("column is present, but is a relationship")
			}

			return attribute.GarTag(), nil
		}
	}

	return "", errors.New("column not present")
}

func (a *CollectionOfAttributes) InnerJoinPresent() bool {
	if len(a.TableChild) > 0 {
		return true
	}

	return false
}

//Remove this in the future
func (a *CollectionOfAttributes) Conditions(values ...interface{}) (conditions []string) {
	var queryValues []interface{}

	for i, condition := range values {
		if i%2 == 0 {
			//conditions = append(conditions, fmt.Sprintf("dmd.dbo.%s.%s", a.Table, condition))
			conditions = append(conditions, fmt.Sprintf("%s.%s", a.Table, condition))
		} else {
			queryValues = append(queryValues, condition)
		}
	}

	return conditions
}

func (a *CollectionOfAttributes) FilterByColumnName(baseName string) (pointerToAttribute interface{}) {
	tableAndColumn := strings.Split(baseName, ".")

	table := tableAndColumn[0]
	columnName := tableAndColumn[1]

	for _, attribute := range a.CollectionOfAttributes {
		if attribute.Table == table {
			if attribute.GarTag() == columnName {
				if addrInterface := attribute.AddrInterface(); addrInterface != nil {
					return addrInterface
				}
			}
		}
	}

	return nil
}
