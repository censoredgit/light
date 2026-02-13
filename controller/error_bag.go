package controller

import "encoding/json"

type ErrorBag struct {
	err        map[string]string
	isModified bool
}

func (t *ErrorBag) Set(name string, err string) {
	t.err[name] = err
	t.isModified = true
}

func (t *ErrorBag) SetRaw(name string, err error) {
	t.err[name] = err.Error()
	t.isModified = true
}

func newErrorBag() *ErrorBag {
	return &ErrorBag{
		err: make(map[string]string),
	}
}

func (t *ErrorBag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.err)
}

func (t *ErrorBag) Errors() map[string]string {
	return t.err
}

func (t *ErrorBag) Has(field string) bool {
	_, ok := t.err[field]
	return ok
}

func (t *ErrorBag) Get(field string) string {
	if t.Has(field) {
		return t.err[field]
	}
	return ""
}

func (t *ErrorBag) SetErrors(errs map[string]string) {
	if errs != nil {
		t.err = errs
		t.isModified = true
	}
}

func (t *ErrorBag) SetRawErrors(errs map[string]error) {
	if errs != nil {
		t.err = make(map[string]string, len(errs))
		for k, v := range errs {
			t.err[k] = v.Error()
		}
		t.isModified = true
	}
}
