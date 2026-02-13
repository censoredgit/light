package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"regexp"
)

type RegexpRule struct {
	Rule
	compiledRegexp *regexp.Regexp
}

func Regexp(pattern string) *RegexpRule {
	return &RegexpRule{
		Rule: Rule{
			alias: "Regexp",
			sType: support.Values,
		},
		compiledRegexp: regexp.MustCompile(pattern),
	}
}

func (r *RegexpRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	for _, value := range inputData.AllValue(fieldName) {
		if !r.compiledRegexp.MatchString(value) {
			return r.Err("The :field field format is invalid.", fieldName)
		}
	}

	return nil
}
