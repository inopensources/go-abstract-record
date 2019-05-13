package composer

import "strings"

type Composer struct {
	Selec      Select
	Update     Update
	Delete     Delete
	From       From
	Set        Set
	Where      Where
	Insert     Insert
	Values     []interface{}
	PostQuery  PostQuery
	PostValues []interface{}
}

func NewComposer() Composer {
	composer := Composer{}
	composer.Selec = NewSelect()
	composer.Insert = NewInsert()
	composer.Update = NewUpdate()
	composer.Delete = NewDelete()

	composer.From = NewFrom()

	composer.Set = NewSet()
	composer.Where = NewWhere()

	return composer
}

func (c *Composer) AddValues(values ...interface{}) {
	c.Values = append(c.Values, values...)
}

func (c *Composer) AddPostValues(values ...interface{}) {
	c.PostValues = append(c.PostValues, values...)
}

func (c *Composer) buildValues() {
	c.Values = append(c.Values, c.PostValues...)
}

func (c *Composer) BuildQuery() (string, []interface{}) {
	var sb strings.Builder

	if c.Selec.Valid() {
		sb.WriteString(c.Selec.Build())
	}

	if c.Update.Valid() {
		sb.WriteString(c.Update.Build())
	}

	if c.Delete.Valid() {
		sb.WriteString(c.Delete.Build())
	}

	if c.From.Valid() {
		sb.WriteString(c.From.Build())
	}

	if c.Set.Valid() {
		sb.WriteString(c.Set.Build())
	}

	if c.Where.Valid() {
		sb.WriteString(c.Where.Build())
	}

	if c.Insert.Valid() {
		sb.WriteString(c.Insert.Build())
	}

	if c.PostQuery.Valid() {
		sb.WriteString(c.PostQuery.Build())
	}

	sb.WriteString(";")

	//Consolidate query values
	c.buildValues()

	return sb.String(), c.Values
}
