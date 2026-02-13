package controller

import (
	"encoding/json"
	"net/url"
)

type InputBag struct {
	data       url.Values
	old        url.Values
	isModified bool
}

func newInputBag() *InputBag {
	return &InputBag{
		data: make(url.Values),
		old:  make(url.Values),
	}
}

func (t *InputBag) MarshalJSON() ([]byte, error) {
	out := make(map[string]map[string][]string)
	out["new"] = make(map[string][]string)
	out["old"] = make(map[string][]string)
	for k, v := range t.data {
		out["new"][k] = v
	}
	for k, v := range t.old {
		out["old"][k] = v
	}
	return json.Marshal(out)
}

func (t *InputBag) All() url.Values {
	return t.data
}

func (t *InputBag) Has(field string) bool {
	_, ok := t.data[field]
	return ok
}

func (t *InputBag) Get(field string) string {
	if data, ok := t.data[field]; ok {
		return data[0]
	}
	return ""
}

func (t *InputBag) List(field string) []string {
	if data, ok := t.data[field]; ok {
		return data
	}
	return []string{}
}

func (t *InputBag) SetOld(old url.Values) {
	t.old = old
	t.isModified = true
}

func (t *InputBag) SetData(data url.Values) {
	t.data = data
	t.isModified = true
}

func (t *InputBag) Set(field, data string) {
	t.data[field] = []string{data}
	t.isModified = true
}

func (t *InputBag) SetList(field string, data []string) {
	t.data[field] = data
	t.isModified = true
}

func (t *InputBag) Old(field string) string {
	return t.OldOrDefault(field, "")
}

func (t *InputBag) OldOrDefault(field, def string) string {
	if old, ok := t.old[field]; ok {
		return old[0]
	}

	if data, ok := t.data[field]; ok {
		return data[0]
	}
	return def
}

func (t *InputBag) OldOrDefaultBool(field string, def bool) bool {
	if old, ok := t.old[field]; ok {
		return old[0] != ""
	}

	if data, ok := t.data[field]; ok {
		return data[0] != ""
	}
	return def
}

func (t *InputBag) OldListIndex(field string, index int) string {
	if old, ok := t.old[field]; ok {
		return old[index]
	}

	if data, ok := t.data[field]; ok {
		return data[index]
	}

	return ""
}

func (t *InputBag) HasEverywhere(field string) bool {
	if _, ok := t.old[field]; ok {
		return true
	}

	if _, ok := t.data[field]; ok {
		return true
	}
	return false
}
