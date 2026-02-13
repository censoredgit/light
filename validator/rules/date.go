package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"time"
)

type DateRule struct {
	Rule
	parameter string
}

func Date(pattern string) *DateRule {
	return &DateRule{
		Rule: Rule{
			alias: "Date",
			sType: support.Values,
		},
		parameter: pattern,
	}
}

func (r *DateRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	for _, date := range inputData.AllValue(fieldName) {
		_, err := time.Parse(r.parameter, date)
		if err != nil {
			return r.Err("The :field field must be a valid date.", fieldName)
		}
	}
	return nil
}
