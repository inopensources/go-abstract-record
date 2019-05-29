package query

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"github.com/infarmasistemas/go-abstract-record/active/query/composer"
	objecto "github.com/infarmasistemas/go-abstract-record/active/query/object"
)

type CompositionOps struct {
	objecto  objecto.Object
	composer composer.Composer
	options  options.OptionsOps
}

func NewCompositionOps(object interface{}, extraOptions ...interface{}) *CompositionOps {
	newCompositionOps := CompositionOps{}
	newCompositionOps.composer = composer.NewComposer()
	newCompositionOps.objecto = objecto.NewObject(object, extraOptions...)
	newCompositionOps.options = options.NewOptionsOps(extraOptions...)

	return &newCompositionOps
}

func (c *CompositionOps) Select(values ...interface{}) (query string, pointerList []interface{}) {
	if c.options.CheckIfCustomFieldsAreFromThisTable(c.objecto.Table) {
		if c.options.QueryCustomFieldsPresent() {
			c.composer.Selec.AddColumn(c.options.QueryCustomFieldsAsSlice(c.objecto.Table)...)
		}
	} else {
		c.composer.Selec.AddColumn(c.objecto.AttributesAsColumnNames()...)
	}

	c.composer.From.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.objecto.Table))

	if len(values) > 0 {
		c.composer.Where.AddCondition(c.objecto.Conditions(values...)...)
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
	c.composer.Insert.AddColumn(c.objecto.AttributesAsColumnNames()...)
	c.composer.Insert.AddTableName(fmt.Sprintf("dmd.dbo.%s ", c.objecto.Table))
	c.composer.Insert.AddValues(c.objecto.AttributeValuesAsArray()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Delete() (query string, pointerList []interface{}) {
	c.composer.Delete.Call()
	c.composer.From.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.objecto.Table))
	c.composer.Where.AddCondition(c.objecto.AttributesAsColumnNames()...)
	c.composer.Where.AddValues(c.objecto.AttributeValuesAsArray()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Update(values ...interface{}) (query string, pointerList []interface{}) {
	c.composer.Update.AddTableName(fmt.Sprintf("dmd.dbo.%s", c.objecto.Table))

	for index, colName := range values {
		if index%2 == 0 {
			c.composer.Set.AddColumn(c.getRealColName(fmt.Sprint(colName)))
		} else {
			c.composer.Set.AddValues(colName)
		}
	}

	c.composer.Where.AddCondition(c.objecto.AttributesAsColumnNames()...)
	c.composer.Where.AddValues(c.objecto.AttributesValues...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) getPointerList() []interface{} {
	return c.objecto.PointerList
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
	return c.objecto.AttributesJson[value]
}
