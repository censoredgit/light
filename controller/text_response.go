package controller

import (
	"net/http"
)

type TextResponse struct {
	data string
	CommonResponse
}

func newTextResponse(data string, code int, flashStorage ContextFlashStorage) *TextResponse {
	return &TextResponse{data: data, CommonResponse: CommonResponse{
		code:         code,
		flashStorage: flashStorage,
	}}
}

func (c *TextResponse) Process(ctx *Ctx) {
	c.CommonResponse.process(ctx)

	ctx.httpResponse.Header().Set("Content-Type", "text/plain; charset=utf-8")

	ctx.httpResponse.WriteHeader(c.Code())
	_, err := ctx.httpResponse.Write([]byte(c.data))
	if err != nil {
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *TextResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}
