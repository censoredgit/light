package rules

import (
	"errors"
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"strings"
)

type Rule struct {
	message string
	alias   string
	sType   support.Type
}

func (r *Rule) SupportType() support.Type {
	return r.sType
}

func (r *Rule) Alias() string {
	return r.alias
}

func (r *Rule) SetMessage(msg string) {
	r.message = msg
}

func (r *Rule) Err(msg string, fieldName string) error {
	if r.message != "" {
		return errors.New(strings.ReplaceAll(r.message, ":field", fieldName))
	}
	return errors.New(strings.ReplaceAll(msg, ":field", fieldName))
}

func (r *Rule) Process(input *input.Data, fieldName string) error {
	return r.Err("The field :field is invalid.", fieldName)
}
