package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
)

type RequiredRule struct {
	*Rule
	allowEmpty bool
}

func Required() *RequiredRule {
	return &RequiredRule{
		Rule: &Rule{
			alias: "Required",
			sType: support.Both,
		},
		allowEmpty: false,
	}
}

func (r *RequiredRule) AllowEmpty() *RequiredRule {
	r.allowEmpty = true
	return r
}

func (r *RequiredRule) Process(inputData *input.Data, fieldName string) error {
	defaultError := r.Err("The :field field is required.", fieldName)

	exists, src := inputData.Has(fieldName)
	if !exists {
		return defaultError
	}

	switch src {
	case input.SourceValues:
		if !inputData.HasValues(fieldName) {
			return defaultError
		}

		if !r.allowEmpty {
			for _, d := range inputData.AllValue(fieldName) {
				if d == "" {
					return defaultError
				}
			}
		}
	default:
		if !inputData.HasFiles(fieldName) {
			return defaultError
		}
	}

	return nil
}
