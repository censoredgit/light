package controller

import (
	"github.com/censoredgit/light/validator"
	"net/http"
)

type RequestValidator struct {
	ruleCollection   *validator.RuleCollection
	jsonResponseCode int
	jsonResponseBody any
	protectFields    []string
}

type ValidatorErrResponse map[string]error

func NewRequestValidator(fn func(ruleCollection *validator.RuleCollection)) *RequestValidator {
	ruleCollection := validator.NewRuleCollection()

	fn(ruleCollection)

	return &RequestValidator{
		ruleCollection:   ruleCollection,
		jsonResponseCode: http.StatusUnauthorized,
		jsonResponseBody: http.StatusText(http.StatusUnauthorized),
	}
}

func (req *RequestValidator) ProtectFields(fields ...string) *RequestValidator {
	req.protectFields = append(req.protectFields, fields...)
	return req
}

func (req *RequestValidator) ProtectedFields() []string {
	return req.protectFields
}

func (req *RequestValidator) RuleCollection() *validator.RuleCollection {
	return req.ruleCollection
}

func (req *RequestValidator) SetJsonResponseCode(code int) *RequestValidator {
	req.jsonResponseCode = code
	return req
}

func (req *RequestValidator) JsonResponseCode() int {
	return req.jsonResponseCode
}

func (req *RequestValidator) SetJsonResponseBody(body any) *RequestValidator {
	req.jsonResponseBody = body
	return req
}

func (req *RequestValidator) JsonResponseBody() any {
	return req.jsonResponseBody
}
