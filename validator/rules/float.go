package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"strconv"
)

type FloatRule struct {
	Rule
}

func Float() *FloatRule {
	return &FloatRule{
		Rule: Rule{
			alias: "Float",
			sType: support.Values,
		},
	}
}

func (r *FloatRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	var err error
	for _, value := range inputData.AllValue(fieldName) {
		if _, err = strconv.ParseFloat(value, 64); err != nil {
			return r.Err("The :field field must be a float.", fieldName)
		}
	}

	return nil
}
