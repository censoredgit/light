package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"

	"slices"
)

type AcceptedRule struct {
	Rule
}

func Accepted() *AcceptedRule {
	return &AcceptedRule{
		Rule: Rule{
			alias: "Accepted",
			sType: support.Values,
		},
	}
}

func (r *AcceptedRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	acceptOptions := []string{"", "1", "on", "true", "yes"}

	for _, val := range inputData.AllValue(fieldName) {
		if slices.Contains(acceptOptions, val) {
			return nil
		}
	}

	return r.Err("The :field field must be accepted.", fieldName)
}
