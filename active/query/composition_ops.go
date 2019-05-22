package query

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/query/composer"
	"reflect"
	"strings"
)

type CompositionOps struct {
	table            string
	attributesGar    []string
	attributesJson   map[string]string
	attributesValues []interface{}
	pointerList      []interface{}
	queryValues      []interface{}
	object           interface{}
	composer         composer.Composer
}

func NewCompositionOps(object interface{}) *CompositionOps {
	newCompositionOps := CompositionOps{}
	newCompositionOps.discoverTable(object)
	newCompositionOps.discoverAttributesAndpointerList(object)
	newCompositionOps.composer = composer.NewComposer()
	newCompositionOps.object = object

	return &newCompositionOps
}

func (c *CompositionOps) Select(values ...interface{}) (query string, pointerList []interface{}) {
	c.composer.Selec.AddColumn(c.attributesAsColumnNames()...)
	c.composer.From.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.table))

	if len(values) > 0 {
		c.composer.Where.AddCondition(c.conditions(values...)...)
		for i := range values {
			if i%2 != 0 {
				//c.composer.AddValues(values[i])
				c.composer.Where.AddValues(values[i])
			}
		}
	}

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Insert() (query string, pointerList []interface{}) {
	c.composer.Insert.AddColumn(c.attributesAsColumnNames()...)
	c.composer.Insert.AddTableName(fmt.Sprintf("dmd.dbo.%s ", c.table))
	c.composer.Insert.AddValues(c.attributeValuesAsArray()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Delete() (query string, pointerList []interface{}) {
	c.composer.Delete.Call()
	c.composer.From.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.table))
	c.composer.Where.AddCondition(c.attributesAsColumnNames()...)
	c.composer.Where.AddValues(c.attributeValuesAsArray()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Update(values ...interface{}) (query string, pointerList []interface{}) {
	c.composer.Update.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.table))

	for index, colName := range values {
		if index%2 == 0 {
			c.composer.Set.AddColumn(c.getRealColName(fmt.Sprint(colName)))
		} else {
			c.composer.Set.AddValues(colName)
		}
	}

	c.composer.Where.AddCondition(c.attributesAsColumnNames()...)
	c.composer.Where.AddValues(c.attributesValues...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) discoverTable(object interface{}) {
	c.table = strings.Split(reflect.TypeOf(object).String(), ".")[len(strings.Split(reflect.TypeOf(object).String(), "."))-1]
}

func (c *CompositionOps) discoverAttributesAndpointerList(object interface{}) {
	var attributeListGar []string

	var fieldList []interface{}
	var attributeValues []interface{}

	c.attributesJson = make(map[string]string, 0)

	s := reflect.ValueOf(object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		finalGar, finalJson := c.parseJsonGar(typeOfT.Field(i))

		if len(finalGar) > 0 && len(finalJson) > 0 {
			someField := s.Field(i)

			fieldList = append(fieldList, someField.Addr().Interface())
			attributeListGar = append(attributeListGar, finalGar)
			c.attributesJson[finalJson] = finalGar

			attributeValues = append(attributeValues, someField.Interface())
		}
	}

	c.attributesGar = attributeListGar
	c.pointerList = fieldList
	c.attributesValues = attributeValues
}

func (c *CompositionOps) parseJsonGar(field reflect.StructField) (string, string) {
	var finalGar, finalJson string

	//Rel is present, therefore it is a relationship model
	if _, presence := field.Tag.Lookup("rel"); presence {
		return finalGar, finalJson
	}

	finalJson = field.Tag.Get("json")
	finalGar = field.Tag.Get("gar")

	return finalGar, finalJson
}

func (c *CompositionOps) attributesAsColumnNames() []string {
	var columns []string
	for _, attributeName := range c.attributesGar {
		columns = append(columns, attributeName)
	}

	return columns
}

func (c *CompositionOps) attributeValuesAsArray() []interface{} {
	var values []interface{}
	for _, value := range c.attributesValues {
		if reflect.ValueOf(value).Elem().IsValid() {
			values = append(values, reflect.ValueOf(value).Elem().Interface())
		} else {
			values = append(values, nil)
		}
	}

	return values
}

func (c *CompositionOps) conditions(values ...interface{}) []string {
	var queryValues []interface{}

	var conditions []string

	for i, condition := range values {
		if i%2 == 0 {
			conditions = append(conditions, fmt.Sprintf("%s", condition))
		} else {
			queryValues = append(queryValues, condition)
		}
	}

	c.queryValues = queryValues
	return conditions
}

func (c *CompositionOps) getPointerList() []interface{} {
	return c.pointerList
}

func (c *CompositionOps) quotedOrNot(value interface{}) string {
	switch v := value.(type) {
	//case int:
	//	fmt.Printf("Twice %v is %v\n", v, v*2)
	//	return fmt.Sprintf("%s", v)
	default:
		//fmt.Printf("I don't know about type %T!\n", v)
		return fmt.Sprintf("'%v'", v)
	}
	return ""
}

func (c *CompositionOps) GetComposer() *composer.Composer {
	return &c.composer
}

func (c *CompositionOps) getRealColName(value string) string {
	return c.attributesJson[value]
}