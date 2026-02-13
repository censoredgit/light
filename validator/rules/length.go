package rules

import (
	"fmt"
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"unicode/utf8"
)

type LengthRule struct {
	Rule
	min     int
	max     int
	asRunes bool
}

func Length(max int) *LengthRule {
	return &LengthRule{
		Rule: Rule{
			alias: "Length",
			sType: support.Values,
		},
		max: max,
	}
}

func (r *LengthRule) SetMin(min int) *LengthRule {
	r.min = min
	return r
}

func (r *LengthRule) AsRunes() *LengthRule {
	r.asRunes = true
	return r
}

func (r *LengthRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	var valueLength int
	for _, value := range inputData.AllValue(fieldName) {
		if r.asRunes {
			valueLength = utf8.RuneCountInString(value)
		} else {
			valueLength = len(value)
		}

		if valueLength > r.max {
			return r.Err(fmt.Sprintf("The :field field must be at least %d.", r.max), fieldName)
		} else if r.min > valueLength {
			return r.Err(fmt.Sprintf("The :field field must not be greater than %d.", r.min), fieldName)
		}
	}

	return nil
}
