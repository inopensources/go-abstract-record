package json_ops

import (
	"encoding/json"
	"fmt"
)

type Parser struct {
	body            string
	parsedBody      map[string]interface{}
	parsedBodyArray []map[string]interface{}
	bodyAsValues    []interface{}
	bodyArrayAsValues [][]interface{}
}

func New(json interface{}) *Parser {
	return &Parser{fmt.Sprintf("%s", json), nil, nil, nil, nil}
}

func (j *Parser) GetBodyAsValues() ([]interface{}, error) {

	err := j.parseBodyToValues()
	if err != nil {
		return nil, err
	}

	return j.bodyAsValues, err
}

func (j *Parser) GetBodyArrayAsValues() ([][]interface{}, error) {

	err := j.parseBodyToValues()
	if err != nil {
		return nil, err
	}

	return j.bodyArrayAsValues, err
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

func (j *Parser) IsArray() bool {
	var js []interface{}
	json.Unmarshal([]byte(j.body), &js)

	return len(js) > 1
}

func (j *Parser) bodyToMap() error {
	var err error
	if j.IsArray() {
		var bodyAsMapArray []map[string]interface{}

		err = json.Unmarshal([]byte(j.body), &bodyAsMapArray)
		if err != nil {
			return err
		}

		j.parsedBodyArray = bodyAsMapArray
	} else {
		var bodyAsMap map[string]interface{}

		err = json.Unmarshal([]byte(j.body), &bodyAsMap)
		if err != nil {
			return err
		}

		j.parsedBody = bodyAsMap
	}

	return err
}

func (j *Parser) mapToValues() {
	if j.IsArray() {
		for _, parsedBody := range j.parsedBodyArray {
			var x []interface{}
			for key, value := range parsedBody {
				x = append(x, key)
				x = append(x, value)
			}

			j.bodyArrayAsValues = append(j.bodyArrayAsValues, x)
		}
	} else {
		for key, value := range j.parsedBody {
			j.bodyAsValues = append(j.bodyAsValues, key)
			j.bodyAsValues = append(j.bodyAsValues, value)
		}
	}
}