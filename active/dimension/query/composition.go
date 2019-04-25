package query

import (
	"fmt"
	"reflect"
	"strings"
)

type Composition struct {
	table       string
	attributes  []string
	pointerList []interface{}
	queryValues []interface{}
}

func NewComposition(object interface{}) *Composition {
	newComposition := Composition{}
	newComposition.discoverTable(object)
	newComposition.discoverAttributesAndpointerList(object)

	return &newComposition
}

func (c *Composition) Select(values ...interface{}) (query string, pointerList []interface{}) {
	var sb strings.Builder

	sb.WriteString("SELECT ")
	sb.WriteString(c.attributesAsSQL())
	sb.WriteString(fmt.Sprintf("FROM dmd.dbo.%s ", c.table))
	sb.WriteString("WHERE ")
	sb.WriteString(c.conditionsAsSQL(values...))
	sb.WriteString(";")

	return sb.String(), c.queryValues
}

func (c *Composition) discoverTable(object interface{}) {
	c.table = strings.Split(reflect.TypeOf(object).String(), ".")[len(strings.Split(reflect.TypeOf(object).String(), "."))-1]
}

func (c *Composition) discoverAttributesAndpointerList(object interface{}) {
	var attributeList []string
	var fieldList []interface{}

	s := reflect.ValueOf(object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		p := strings.Replace(string(typeOfT.Field(i).Tag), "json:\"", "", -1)
		p = strings.Replace(p, "\"", "", -1)

		someField := s.Field(i)
		fieldList = append(fieldList, someField.Addr().Interface())
		attributeList = append(attributeList, p)
	}

	c.attributes = attributeList
	c.pointerList = fieldList
}

func (c *Composition) attributesAsSQL() string {
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

func (c *Composition) conditionsAsSQL(values ...interface{}) string {
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

func (c *Composition) getPointerList() []interface{} {
	return c.pointerList
}