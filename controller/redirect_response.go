package controller

import (
	"net/http"
)

type RedirectResponse struct {
	url string
	CommonResponse
}

func newRedirect(url string, flashStorage ContextFlashStorage) *RedirectResponse {
	return &RedirectResponse{url: url, CommonResponse: CommonResponse{
		code:         http.StatusFound,
		flashStorage: flashStorage,
	}}
}

func (c *RedirectResponse) Process(ctx *Ctx) {
	c.CommonResponse.process(ctx)

	http.Redirect(ctx.httpResponse, ctx.request, c.url, c.code)
}

func (c *RedirectResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}
