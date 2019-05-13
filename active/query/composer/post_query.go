package composer

import "strings"

type PostQuery struct {
	body 	[]string
}

func NewPostQuery() PostQuery {
	postQuery := PostQuery{}

	return postQuery
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
