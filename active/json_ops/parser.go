package json_ops

import (
	"encoding/json"
	"fmt"
)

type Parser struct {
	body         string
	parsedBody   map[string]interface{}
	bodyAsValues []interface{}
}

func New(json interface{}) *Parser {
	return &Parser{fmt.Sprintf("%s", json), nil, nil}
}

func (j *Parser) GetBodyAsValues() ([]interface{}, error) {

	err := j.parseBodyToValues()
	if err != nil {
		return nil, err
	}

	return j.bodyAsValues, err
}

func (j *Parser) parseBodyToValues() error {

	err := j.bodyToMap()
	if err != nil {
		return err
	}

	j.mapToValues()

	return err
}

func (j *Parser) IsJSON() bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(j.body), &js) == nil
}

func (j *Parser) bodyToMap() error {
	bodyAsMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(j.body), &bodyAsMap)
	if err != nil {
		return err
	}

	j.parsedBody = bodyAsMap

	return err
}

func (j *Parser) mapToValues() {
	for key, value := range j.parsedBody {
		j.bodyAsValues = append(j.bodyAsValues, key)
		j.bodyAsValues = append(j.bodyAsValues, value)
	}
}
