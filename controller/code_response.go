package controller

import (
	"fmt"
	"net/http"
)

type CodeResponse struct {
	CommonResponse
}

func (c *CodeResponse) Process(ctx *Ctx) {
	c.CommonResponse.process(ctx)

	ctx.httpResponse.WriteHeader(c.code)
	_, _ = fmt.Fprintln(ctx.httpResponse, http.StatusText(c.code))
}

func newCodeResponse(code int, flashStorage ContextFlashStorage) *CodeResponse {
	resp := &CodeResponse{
		CommonResponse{
			flashStorage: flashStorage,
		},
	}
	resp.code = code
	return resp
}

func (c *CodeResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}
