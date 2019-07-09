package check

import (
	"errors"
	"fmt"
	"github.com/infarmasistemas/go-abstract-record/active/check/json_ops"
)

type CheckOps struct {
}

func (a *CheckOps) CheckAndExecute(function func(finalValues ...interface{}) error, values ...interface{}) error {

	// If len == 1, then it's either JSON or pure SQL
	if len(values) == 1 {
		// If a JSON is coming in...
		if json_ops.New(values[0]).IsJSON() {
			//Check if JSON contains more than one collection_of_attributes

			finalValues, err := json_ops.New(values[0]).GetBodyAsValues()
			if err != nil {
				return err
			}

			return function(finalValues...)
		}

		// If it's not a JSON, then it's pure sql
		// TODO: Do something here
		return errors.New(fmt.Sprintf("Can't deal with pure sql %s", values))
	}

	return function(values...)
}

func (a *CheckOps) TreatEntry(functionSingle func(json []byte) error, functionMultiple func(json []byte) error, values ...interface{}) error {
	if len(values) == 1 {
		if json_ops.New(values[0]).IsJSON() {

			if json_ops.New(values[0]).IsArray() {
				//JSON is an array, therefore I need to write data to an array
				return functionMultiple([]byte(fmt.Sprintf("%s", values[0])))
			}

			return functionSingle([]byte(fmt.Sprintf("%s", values[0])))
		}

		return errors.New(fmt.Sprintf("Can't deal with pure sql %s", values))
	}

	return functionSingle([]byte(fmt.Sprintf("%s", values[0])))
}

func (a *CheckOps) TreatValuesForUpdate(fn func(fnValues ...interface{}) error, values ...interface{}) error {
	vls, err := json_ops.New(values[0]).GetBodyAsValues()
	if err != nil {
		return err
	}

	return fn(vls...)
}
