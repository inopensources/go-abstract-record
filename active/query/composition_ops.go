package query

import (
	"fmt"
	"reflect"
	"strings"
)

type CompositionOps struct {
	table            string
	attributes       []string
	attributesValues []interface{}
	pointerList      []interface{}
	queryValues      []interface{}
	queriesOps       QueriesOps
}

func NewCompositionOps(object interface{}) *CompositionOps {
	newCompositionOps := CompositionOps{}
	newCompositionOps.discoverTable(object)
	newCompositionOps.discoverAttributesAndpointerList(object)

	return &newCompositionOps
}

func (c *CompositionOps) Select(values ...interface{}) (query string, pointerList []interface{}) {
	c.queriesOps.AddPreQuery("SELECT ")
	c.queriesOps.AddPreQuery(c.attributesAsSQL())
	c.queriesOps.AddPreQuery(fmt.Sprintf("FROM dmd.dbo.%s ", c.table))

	if len(values) > 0 {
		c.queriesOps.AddMidQuery("WHERE ")
		c.queriesOps.AddMidQuery(c.conditionsAsSQL(values...))
		c.queriesOps.AddValues(values...)
	}

	return c.queriesOps.returnBuiltQueryAndValues()
}

func (c *CompositionOps) Insert() (query string, pointerList []interface{}) {
	c.queriesOps.AddPreQuery("INSERT INTO ")
	c.queriesOps.AddPreQuery(fmt.Sprintf("dmd.dbo.%s ", c.table))
	c.queriesOps.AddPreQuery("(")
	c.queriesOps.AddPreQuery(c.attributesAsSQL())
	c.queriesOps.AddPreQuery(") ")

	c.queriesOps.AddMidQuery("VALUES ")
	c.queriesOps.AddMidQuery("(")
	c.queriesOps.AddMidQuery(c.attributeValuesAsSQL())
	c.queriesOps.AddMidQuery(")")

	c.queriesOps.AddValues(c.attributeValuesAsArray()...)

	return c.queriesOps.returnBuiltQueryAndValues()
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

func (c *CompositionOps) attributesAsSQL() string {
	var sb strings.Builder
	for i, x := range c.attributes {
		sb.WriteString(fmt.Sprintf("%s", x))
		if i < len(c.attributes)-1 {
			sb.WriteString(", ")
		} else {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}

func (c *CompositionOps) attributeValuesAsSQL() string {
	var sb strings.Builder
	for i, _ := range c.attributesValues {
		sb.WriteString(fmt.Sprintf("%s","?"))
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
		values = append(values, value)
	}

	return values
}

func (c *CompositionOps) conditionsAsSQL(values ...interface{}) string {
	var sb strings.Builder
	var queryValues []interface{}

	for i, x := range values {
		if i%2 == 0 {
			sb.WriteString(fmt.Sprintf("%s = ?", x))
			if i+1 < len(values)-1 {
				sb.WriteString(" AND ")
			}
		} else {
			queryValues = append(queryValues, x)
		}
	}

	c.queryValues = queryValues
	return sb.String()
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

func (c *CompositionOps) GetQueriesOps() *QueriesOps {
	return &c.queriesOps
}
