package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"strconv"
)

type IntegerRule struct {
	Rule
}

func Integer() *IntegerRule {
	return &IntegerRule{
		Rule: Rule{
			alias: "Integer",
			sType: support.Values,
		},
	}
}

func (r *IntegerRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	for _, integer := range inputData.AllValue(fieldName) {
		_, err := strconv.Atoi(integer)
		if err != nil {
			return r.Err("The :field field must be an integer.", fieldName)
		}
	}

	return nil
}
