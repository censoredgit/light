package controller

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	data     interface{}
	errorBag *ErrorBag
	inputBag *InputBag
	CommonResponse
}

func (c *JsonResponse) Errors() *ErrorBag {
	return c.errorBag
}

func (c *JsonResponse) Inputs() *InputBag {
	return c.inputBag
}

func (c *JsonResponse) Flush() {

}

func newJsonResponse(data interface{}, code int) *JsonResponse {
	response := &JsonResponse{
		data:     data,
		inputBag: newInputBag(),
		errorBag: newErrorBag(),
		CommonResponse: CommonResponse{
			code: code,
		}}

	response.flashStorage = response

	return response
}

func (c *JsonResponse) Process(ctx *Ctx) {
	c.CommonResponse.process(ctx)

	ctx.httpResponse.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, err := json.Marshal(struct {
		Data   interface{} `json:"data"`
		Errors *ErrorBag   `json:"errors"`
		Input  *InputBag   `json:"input"`
	}{
		Data:   c.data,
		Errors: c.errorBag,
		Input:  c.inputBag,
	})
	if err != nil {
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.httpResponse.WriteHeader(c.Code())
	_, err = ctx.httpResponse.Write(b)
	if err != nil {
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *JsonResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}
