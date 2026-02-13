package controller

import (
	"net/http"
)

type HtmlResponse struct {
	data string
	CommonResponse
}

func newHtmlResponse(data string, code int, flashStorage ContextFlashStorage) *HtmlResponse {
	return &HtmlResponse{data: data, CommonResponse: CommonResponse{
		code:         code,
		flashStorage: flashStorage,
	}}
}

func (c *HtmlResponse) Process(ctx *Ctx) {
	c.CommonResponse.process(ctx)

	ctx.httpResponse.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx.httpResponse.WriteHeader(c.Code())
	_, err := ctx.httpResponse.Write([]byte(c.data))
	if err != nil {
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *HtmlResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}
