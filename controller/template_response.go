package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/flosch/pongo2/v6"
)

type TemplateResponse struct {
	view string
	*CommonResponse
}

func newTemplateResponse(view string, flashStorage ContextFlashStorage) *TemplateResponse {
	return &TemplateResponse{view: view, CommonResponse: &CommonResponse{
		code:         http.StatusOK,
		flashStorage: flashStorage,
	}}
}

func (c *TemplateResponse) View() string {
	return c.view
}

func (c *TemplateResponse) Process(ctx *Ctx) {
	c.CommonResponse.process(ctx)

	if templateSet == nil {
		ctx.Log().Error("should setup template path")
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
		return
	}

	tpl, err := templateSet.FromFile(c.view)

	if err != nil {
		ctx.Log().Error(err.Error())
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Exec(ctx, tpl)
}

func (c *TemplateResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}

func (c *TemplateResponse) Exec(ctx *Ctx, tpl *pongo2.Template) {
	pctx := pongo2.Context{
		"ctx":     ctx,
		"Data":    c.Data(),
		"Err":     ctx.flashStorage.Errors(),
		"Input":   ctx.flashStorage.Inputs(),
		"IsAuth":  func() bool { return ctx.IsAuth() },
		"IsGuest": func() bool { return ctx.IsGuest() },
		"Role": func(role string) bool {
			return ctx.Role(role)
		},
		"Can": func(permission string) bool { return ctx.Can(permission) },
		"StaticFile": func(targetPath string) string {
			return path.Join(config.StaticPath, targetPath)
		},
		"Route": func(name string, args ...interface{}) string {
			return ctx.Route(name, args...)
		},
		"Url": func(routeName string, args ...interface{}) string {
			return ctx.Url(routeName, args...)
		},
		"HasCurrentRoutePrefix": func(name string, args ...interface{}) bool {
			return ctx.HasCurrentRoutePrefix(name, args...)
		},
		"IsCurrentRoute": func(name string) bool {
			return ctx.IsCurrentRoute(name)
		},
		"AuthIdentification": func() string {
			return ctx.AuthIdentification()
		},
		"CsrfToken": func() string {
			return ctx.CsrfToken()
		},
		"CsrfTokenQuery": func() string {
			return url.QueryEscape(ctx.CsrfFieldName()) + "=" + ctx.CsrfToken()
		},
		"CsrfTokenInput": func() *pongo2.Value {
			return pongo2.AsSafeValue(
				fmt.Sprint("<input type=\"hidden\" name=\"", ctx.CsrfFieldName(), "\" value=\"", ctx.CsrfToken(), "\" />"))
		},
	}

	for k, v := range templateFuncMap {
		pctx[k] = v
	}

	ctx.httpResponse.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.httpResponse.WriteHeader(c.Code())

	err := tpl.ExecuteWriter(
		pctx,
		ctx.httpResponse,
	)

	if err != nil {
		ctx.Log().Error(err.Error())
		ctx.httpResponse.WriteHeader(http.StatusInternalServerError)
	}
}
