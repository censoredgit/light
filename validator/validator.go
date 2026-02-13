package validator

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Rule interface {
	Alias() string
	Process(inputData *input.Data, fieldName string) error
	SetMessage(msg string)
	SupportType() support.Type
}

type Validator struct {
	isFailed       bool
	errs           map[string]string
	ruleCollection *RuleCollection
	inpData        *input.Data
}

func New() *Validator {
	return &Validator{ruleCollection: NewRuleCollection(),
		errs: make(map[string]string)}
}

func (v *Validator) AddRule(field string, rules ...Rule) *Validator {
	v.ruleCollection.AddRule(field, rules...)
	return v
}

func (v *Validator) SetRuleCollection(ruleCollection *RuleCollection) *Validator {
	v.ruleCollection = ruleCollection
	return v
}

func (v *Validator) Validate(inpData *input.Data) bool {
	v.inpData = inpData

	var err error
	for field, rules := range *v.ruleCollection {
		for _, rule := range rules {
			err = rule.Process(inpData, field)
			if err != nil {
				v.errs[field] = err.Error()
				v.isFailed = true
				break
			}
		}
	}
	return !v.isFailed
}

func (v *Validator) ValidateByRequestForms(formValues *url.Values, multipartForm *multipart.Form) bool {
	return v.Validate(input.FormsToInputData(formValues, multipartForm))
}

func (v *Validator) ValidateByJsonRequest(r *http.Request) bool {
	return v.Validate(input.RequestToInputData(r))
}

func (v *Validator) Errors() map[string]string {
	return v.errs
}

func (v *Validator) IsFailed() bool {
	return v.isFailed
}

func (v *Validator) Fields() []string {
	fields := make([]string, len(*v.ruleCollection))
	index := 0
	for field := range *v.ruleCollection {
		fields[index] = field
		index++
	}
	return fields
}

func (v *Validator) InputData() *input.Data {
	return v.inpData
}
