package input

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

const (
	SourceValues = "values"
	SourceFiles  = "files"
)

type Data struct {
	values map[string][]string
	files  map[string][]*multipart.FileHeader
}

func NewInputData() *Data {
	return &Data{
		values: make(map[string][]string),
		files:  make(map[string][]*multipart.FileHeader),
	}
}

func (d *Data) HasString(s string) bool {
	_, ok := d.values[s]
	return ok
}

func (d *Data) HasFile(s string) bool {
	_, ok := d.files[s]
	return ok
}

func (d *Data) GetString(s string) string {
	return d.values[s][0]
}

func (d *Data) GetInt(s string) (int, error) {
	return strconv.Atoi(d.GetString(s))
}

func (d *Data) GetFile(s string) *multipart.FileHeader {
	return d.files[s][0]
}

func (d *Data) GetStrings(s string) []string {
	return d.values[s]
}

func (d *Data) GetFiles(s string) []*multipart.FileHeader {
	return d.files[s]
}

func (d *Data) Has(field string) (bool, string) {
	var source string
	_, ok := d.files[field]
	source = SourceFiles
	if !ok {
		_, ok = d.values[field]
		source = SourceValues
	}
	return ok, source
}

func (d *Data) HasValues(field string) bool {
	_, ok := d.values[field]
	return ok
}

func (d *Data) HasFiles(field string) bool {
	_, ok := d.files[field]
	return ok
}

func (d *Data) AllValue(field string) []string {
	return d.values[field]
}

func (d *Data) Value(field string) string {
	return d.values[field][0]
}

func (d *Data) AllFiles(field string) []*multipart.FileHeader {
	return d.files[field]
}

func (d *Data) File(field string) *multipart.FileHeader {
	return d.files[field][0]
}

func (d *Data) SetValue(field string, data ...string) {
	d.values[field] = data
}

func (d *Data) SetFile(field string, data ...*multipart.FileHeader) {
	d.files[field] = data
}

func FormsToInputData(formValues *url.Values, multipartForm *multipart.Form) *Data {
	inputData := NewInputData()

	for field, val := range *formValues {
		inputData.values[field] = val
	}

	if multipartForm != nil {
		for field, val := range multipartForm.File {
			inputData.files[field] = val
		}
	}

	return inputData
}

func RequestToInputData(r *http.Request) *Data {
	inputData := NewInputData()

	result := make(map[string]any)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&result)
	if err != nil {
		return inputData
	}

	for key, val := range result {
		inputData.SetValue(key, fmt.Sprintf("%v", val))
	}

	err = r.ParseForm()
	if err == nil {
		for k, v := range r.Form {
			inputData.SetValue(k, v...)
		}
	}

	return inputData
}
