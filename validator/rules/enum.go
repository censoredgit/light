package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"slices"
)

type EnumRule struct {
	Rule
	enum []string
}

func Enum(enum ...string) *EnumRule {
	return &EnumRule{
		Rule: Rule{
			alias: "Enum",
			sType: support.Values,
		},
		enum: enum,
	}
}

func (r *EnumRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	for _, enum := range inputData.AllValue(fieldName) {
		if slices.Contains(r.enum, enum) {
			return nil
		}
	}

	return r.Err("The selected :field is invalid.", fieldName)
}
