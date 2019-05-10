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
	extraFuncs  	 []func() string
}

func NewCompositionOps(object interface{}, extraFuncs ...func() string) *CompositionOps {
	newCompositionOps := CompositionOps{}
	newCompositionOps.discoverTable(object)
	newCompositionOps.discoverAttributesAndpointerList(object)
	newCompositionOps.extraFuncs = extraFuncs

	return &newCompositionOps
}

func (c *CompositionOps) Select(values ...interface{}) (query string, pointerList []interface{}) {
	var sb strings.Builder

	sb.WriteString("SELECT ")
	sb.WriteString(c.attributesAsSQL())
	sb.WriteString(fmt.Sprintf("FROM dmd.dbo.%s ", c.table))

	if len(values) > 0 {
		sb.WriteString("WHERE ")
		sb.WriteString(c.conditionsAsSQL(values...))
	}

	if len(c.extraFuncs) > 0 {
		sb.WriteString(" ")
		for _, function := range c.extraFuncs {
			sb.WriteString(function())
		}
	}

	sb.WriteString(";")

	return sb.String(), c.queryValues
}

func (c *CompositionOps) Insert() (query string, pointerList []interface{}) {
	var sb strings.Builder

	sb.WriteString("INSERT INTO ")
	sb.WriteString(fmt.Sprintf("dmd.dbo.%s ", c.table))
	sb.WriteString("(" + c.attributesAsSQL() + ") ")
	sb.WriteString("VALUES ")
	sb.WriteString("(" + c.attributeValuesAsSQL() + ")")
	sb.WriteString(";")

	return sb.String(), c.queryValues
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
	for i, x := range c.attributesValues {
		sb.WriteString(fmt.Sprintf("%s", c.quotedOrNot(x)))
		if i < len(c.attributes)-1 {
			sb.WriteString(", ")
		} else {
			sb.WriteString(" ")
		}
	}

	return sb.String()
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
