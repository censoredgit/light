package rules

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"net/mail"
)

type EmailRule struct {
	Rule
}

func Email() *EmailRule {
	return &EmailRule{
		Rule: Rule{
			alias: "Email",
			sType: support.Values,
		},
	}
}

func (r *EmailRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceFiles {
		return r.Err("The :field field is invalid.", fieldName)
	}

	for _, email := range inputData.AllValue(fieldName) {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return r.Err("The :field field must be a valid email address.", fieldName)
		}
	}

	return nil
}
