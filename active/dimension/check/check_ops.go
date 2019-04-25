package check

import (
	"errors"
	"fmt"
	"github.com/infarmasistemas/goorm/active/json_ops"
)

type CheckOps struct {
}

func (a *CheckOps) CheckAndExecute(function func(finalValues ...interface{}) error, values ...interface{}) error {

	// If len == 1, then it's either JSON or pure SQL
	if len(values) == 1 {
		// If a JSON is coming in...
		if json_ops.New(values[0]).IsJSON() {
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
