package query

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/query/composer"
	"reflect"
	"strings"
)

type CompositionOps struct {
	table            string
	attributes       []string
	attributesValues []interface{}
	pointerList      []interface{}
	queryValues      []interface{}
	composer         composer.Composer
}

func NewCompositionOps(object interface{}) *CompositionOps {
	newCompositionOps := CompositionOps{}
	newCompositionOps.discoverTable(object)
	newCompositionOps.discoverAttributesAndpointerList(object)
	newCompositionOps.composer = composer.NewComposer()

	return &newCompositionOps
}

func (c *CompositionOps) Select(values ...interface{}) (query string, pointerList []interface{}) {
	c.composer.Selec.AddColumn(c.attributesAsColumnNames()...)
	c.composer.From.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.table))

	if len(values) > 0 {
		fmt.Println(c.conditions())
		c.composer.Where.AddCondition(c.conditions(values...)...)
		for i := range values {
			if i % 2 != 0 {
				c.composer.AddValues(values[i])
			}
		}
	}

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Insert() (query string, pointerList []interface{}) {
	c.composer.Insert.AddColumn(c.attributesAsColumnNames()...)
	c.composer.Insert.AddTableName(fmt.Sprintf("dmd.dbo.%s ", c.table))
	c.composer.AddValues(c.attributeValuesAsArray()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Delete() (query string, pointerList []interface{}) {
	c.composer.Delete.Call()
	c.composer.From.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.table))
	c.composer.Where.AddCondition(c.attributesAsColumnNames()...)
	c.composer.AddValues(c.attributeValuesAsArray()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Update() (query string, pointerList []interface{}) {
	c.composer.Update.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.table))
	c.composer.Set.AddColumn(c.attributesAsColumnNames()...)
	c.composer.Where.AddCondition(c.attributesAsColumnNames()...)

	//Review this
	c.composer.AddValues(c.attributeValuesAsArray()...)
	c.composer.AddValues(c.attributeValuesAsArray()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) discoverTable(object interface{}) {
	c.table = strings.Split(reflect.TypeOf(object).String(), ".")[len(strings.Split(reflect.TypeOf(object).String(), "."))-1]
}

func (c *CompositionOps) discoverAttributesAndpointerList(object interface{}) {
	var attributeList []string
	var fieldList []interface{}
	var attributeValues []interface{}

	s := reflect.ValueOf(object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		p := strings.Replace(string(typeOfT.Field(i).Tag), "json:\"", "", -1)
		p = strings.Replace(p, "\"", "", -1)

		someField := s.Field(i)
		fieldList = append(fieldList, someField.Addr().Interface())
		attributeList = append(attributeList, p)
		attributeValues = append(attributeValues, someField.Interface())
	}

	c.attributes = attributeList
	c.pointerList = fieldList
	c.attributesValues = attributeValues
}

func (c *CompositionOps) attributesAsColumnNames() []string {
	var columns []string
	for _, attributeName := range c.attributes {
		columns = append(columns, attributeName)
	}

	return columns
}

func (c *CompositionOps) attributeValuesAsSQL() string {
	var sb strings.Builder
	for i, _ := range c.attributesValues {
		sb.WriteString(fmt.Sprintf("%s", "?"))
		if i < len(c.attributes)-1 {
			sb.WriteString(", ")
		} else {
			sb.WriteString(" ")
		}
	}

	return sb.String()
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

func (c *CompositionOps) GetComposer() *composer.Composer{
	return &c.composer
}
