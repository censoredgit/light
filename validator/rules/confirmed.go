package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
)

type ConfirmedRule struct {
	Rule
	field string
}

func Confirmed(field string) *ConfirmedRule {
	return &ConfirmedRule{
		Rule: Rule{
			alias: "Confirmed",
			sType: support.Values,
		},
		field: field,
	}
}

func (r *ConfirmedRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	defaultError := r.Err("The :field field confirmation does not match.", fieldName)

	if ok := inputData.HasValues(r.field); !ok {
		return defaultError
	}

	data := inputData.AllValue(fieldName)
	confirmedValues := inputData.AllValue(r.field)

	if len(data) != len(confirmedValues) {
		return defaultError
	}

	for index, v := range data {
		if confirmedValues[index] != v {
			return defaultError
		}
	}

	return nil
}
