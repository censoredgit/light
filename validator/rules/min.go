package rules

import (
	"fmt"
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"strconv"
)

type MinRule struct {
	Rule
	parameter float64
}

func Min(min float64) *MinRule {
	return &MinRule{
		Rule: Rule{
			alias: "Min",
			sType: support.Both,
		},
		parameter: min,
	}
}

func (r *MinRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	switch src {
	case input.SourceValues:
		if !inputData.HasValues(fieldName) {
			return nil
		}

		for _, value := range inputData.AllValue(fieldName) {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return r.Err("The :field field must be a number.", fieldName)
			}
			if f < r.parameter {
				return r.Err(fmt.Sprintf("The :field field must be at least %.0f.", r.parameter), fieldName)
			}
		}
	default:
		if !inputData.HasFiles(fieldName) {
			return nil
		}

		for _, value := range inputData.AllFiles(fieldName) {
			if value.Size < int64(r.parameter) {
				return r.Err(fmt.Sprintf("The :field field must be at least %.0f.", r.parameter), fieldName)
			}
		}
	}

	return nil
}
