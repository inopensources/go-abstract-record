package composer

import (
	"github.com/infarmasistemas/go-abstract-record/active/query/composer/object_value"
	"strings"
)

type PostQuery struct {
	body 	[]string
	objectValues []object_value.ObjectValue
}

func NewPostQuery() PostQuery {
	postQuery := PostQuery{}

	return postQuery
}

func (p *PostQuery) AddValues(values ...interface{}) {
	for _, value := range values {
		p.objectValues = append(p.objectValues, object_value.NewObjectValue(value))
	}
}

func (p *PostQuery) AddPostQuery(value ...string) {
	p.body = append(p.body, value...)
}

func (p *PostQuery) Valid() bool {
	if len(p.body) > 0 {
		return true
	} else {
		return false
	}
}

func (p *PostQuery) Build() string {
	var sb strings.Builder

	// Writing body
	for _, body := range p.body {
		sb.WriteString(body)
	}

	return sb.String()
}

func (p *PostQuery) getValues() []interface{} {
	var values []interface{}
	for _, value := range p.objectValues {
		values = append(values, value.GetObject())

	}

	return values
}
