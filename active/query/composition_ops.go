package query

import (
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/options"
	"github.com/infarmasistemas/go-abstract-record/active/query/collection_of_attributes"
	"github.com/infarmasistemas/go-abstract-record/active/query/composer"
)

type CompositionOps struct {
	CollectionOfAttributes collection_of_attributes.CollectionOfAttributes
	composer               composer.Composer
	options                options.OptionsOps
}

func NewCompositionOps(object interface{}, extraOptions ...interface{}) *CompositionOps {
	newCompositionOps := CompositionOps{}
	newCompositionOps.composer = composer.NewComposer()
	newCompositionOps.CollectionOfAttributes = collection_of_attributes.NewCollectionOfAttributes(object, extraOptions...)
	newCompositionOps.options = options.NewOptionsOps(extraOptions...)

	return &newCompositionOps
}

func (c *CompositionOps) Select(values ...interface{}) (query string, pointerList []interface{}) {
	if c.options.CheckIfCustomFieldsAreFromThisTable(c.CollectionOfAttributes.Table) {
		if c.options.QueryCustomFieldsPresent() {
			c.composer.Selec.AddColumn(c.options.QueryCustomFieldsAsSlice(c.CollectionOfAttributes.Table)...)
		}
	} else {
		c.composer.Selec.AddColumn(c.CollectionOfAttributes.AttributesAsColumnNamesForSelect()...)
	}

	c.composer.From.AddTableName(fmt.Sprintf("%s", c.CollectionOfAttributes.Table))

	if c.CollectionOfAttributes.InnerJoinPresent() {
		c.composer.Join.SetParentTable(c.CollectionOfAttributes.Table)
		c.composer.Join.SetChildTables(c.CollectionOfAttributes.TableChild...)
	}

	if len(values) > 0 {
		c.composer.Where.AddCondition(c.CollectionOfAttributes.Conditions(values...)...)
		for i := range values {
			if i%2 != 0 {
				c.composer.Where.AddValues(values[i])
			}
		}
	}

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Count(values ...interface{}) (query string, pointerList []interface{}) {
	c.composer.Count.AddColumn("*")

	c.composer.From.AddTableName(fmt.Sprintf("%s", c.CollectionOfAttributes.Table))

	if len(values) > 0 {
		c.composer.Where.AddCondition(c.CollectionOfAttributes.Conditions(values...)...)
		for i := range values {
			if i%2 != 0 {
				c.composer.Where.AddValues(values[i])
			}
		}
	}

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Insert() (query string, pointerList []interface{}) {
	c.composer.Insert.AddColumn(c.CollectionOfAttributes.AttributesAsColumnNamesForInsert()...)
	c.composer.Insert.AddTableName(fmt.Sprintf("%s ", c.CollectionOfAttributes.Table))

	for _, colName := range c.CollectionOfAttributes.AttributesAsColumnNamesForInsert() {
		c.composer.Insert.AddValues(c.CollectionOfAttributes.AttributeValueFromColumnName(colName))
	}

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Delete() (query string, pointerList []interface{}) {
	c.composer.Delete.Call()
	c.composer.From.AddTableName(fmt.Sprintf("%s", c.CollectionOfAttributes.Table))

	c.composer.Where.AddCondition(c.CollectionOfAttributes.AttributesAsColumnNamesForInsert()...)

	for _, colName := range c.CollectionOfAttributes.AttributesAsColumnNamesForInsert() {
		c.composer.Where.AddValues(c.CollectionOfAttributes.AttributeValueFromColumnName(colName))
	}

	return c.composer.BuildQuery()
}

func (c *CompositionOps) Update(values ...interface{}) (query string, pointerList []interface{}) {
	c.composer.Update.AddTableName(fmt.Sprintf("%s", c.CollectionOfAttributes.Table))

	for index := 0; index < len(values); index += 2 {
		colName := fmt.Sprint(values[index])

		if colName, err := c.CollectionOfAttributes.RealColName(colName); err == nil {
			c.composer.Set.AddColumn(colName)

			passedValue := values[index+1]
			c.composer.Set.AddValues(passedValue)
		}
	}

	c.composer.Where.AddCondition(c.CollectionOfAttributes.AttributesAsColumnNamesForInsert()...)
	c.composer.Where.AddValues(c.CollectionOfAttributes.AttributeValuesAsSlice()...)

	return c.composer.BuildQuery()
}

func (c *CompositionOps) GetComposer() *composer.Composer {
	return &c.composer
}
