[![Go Report Card](https://goreportcard.com/badge/github.com/infarmasistemas/go-abstract-record)](https://goreportcard.com/report/github.com/infarmasistemas/go-abstract-record)

# ORM-like solution to work with legacy SQL databases

This is our in-house ORM-like WIP solution to integrate Go and our legacy SQL database.

The challenge of developing new applications using newer languages like GO without having 
the option of changing/replacing the database structure currently being used is a hard one to tackle.

Thus we decided to write a solution that would:
 - Work as a middle layer between our legacy database and our Go code
 - Integrate seamlessly with Iris
 - Allow us to map our struct models to whatever table structure regardless of the 
   column names being used
 - Let us work with nil values without using SQL's null data types
 - Allow us to set default attribute values via tags regardless of the attribute type
 - Let us render models with their relationships as composed attributes
 - Offer us basic ActiveRecord-like methods like Find, All, Count, Max, Save, Delete, Update...

Basic usage example

Your model should look like this:

````
package models

import (
	database "api-atacado/datasource"
	"fmt"
	"time"
	"github.com/infarmasistemas/go-abstract-record/active"
)

type Clien struct {
    //PS: You should always use pointers to values
    
    //Notice that the gar tag holds the name of the column that you want to map
    SomeCode    *int       `json:"SomeCode" gar:"Some_Code"`
    
    //You can also set default values to attributes by using the default_value tag
    SomeString  *string    `json:"SomeString" gar:"Some_String" default_value:"SomeStringHERE"`
}

func (c *Clien) prepare(objectArray interface{}, extraOptions ...interface{}) *active.AbstractRecord {
	var clients []Clien
	return active.NewAbstractRecord(c, clients, databaseConnection, extraOptions...)
}

func (c *Clien) All(offset, pageSize int, values ...interface{}) ([]Clien, error) {
	var clients []Clien
	return clients, c.prepare(&clients, true, 1).Inner().Paginate("Clien.Some_Code", offset, pageSize).Where(values...)
}

func (c *Clien) Count(values ...interface{}) (int, error) {
	return c.prepare(nil).Count(values...)
}

func (c *Clien) Max(value interface{}) (int, error) {
	return c.prepare(nil).Max(value)
}

func (c *Clien) NewSingle(values ...interface{}) (Clien, error) {
	if len(values) == 0 {
		return *c, nil
	}

	return *c, c.prepare(nil).New(values...)
}

func (c *Clien) NewWithDefaults() (Clien, error) {
	return *c, c.prepare(nil).NewWithDefaults()
}

func (c *Clien) Find(values ...interface{}) (Clien, error) {
	return *c, c.prepare(nil, true, 1).Inner().Find(values...)
}

func (c *Clien) Save() error {
	return c.prepare(nil).Save()
}

func (c *Clien) Delete() error {
	return c.prepare(nil).Delete()
}

func (c *Clien) Update(values ...interface{}) error {
	return c.prepare(nil).Update(values...)
}

````

 
Criação de um framework embrionário da Infarma Sistemas e Serviços
