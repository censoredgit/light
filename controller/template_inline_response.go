package controller

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
)

type TemplateInlineResponse struct {
	view string
	*TemplateResponse
}

func newTemplateInlineResponse(view string, flashStorage ContextFlashStorage) *TemplateInlineResponse {
	return &TemplateInlineResponse{view: view, TemplateResponse: &TemplateResponse{
		CommonResponse: &CommonResponse{
			code:         http.StatusOK,
			flashStorage: flashStorage,
		},
	}}
}

func (c *TemplateInlineResponse) View() string {
	return c.view
}

func (c *TemplateInlineResponse) Process(ctx *Ctx) {
	c.CommonResponse.process(ctx)

	tpl, err := pongo2.FromString(c.view)

	if err != nil {
		ctx.Log().Error(err.Error())
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Exec(ctx, tpl)
}

func (c *TemplateInlineResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}
