package controller

import (
	"github.com/censoredgit/light/utils"
	"net/http"
	"strings"
)

var csrfMiddlewareCfg struct {
	tokenFieldName string
}

func setupCsrfMiddleware(tokenFieldName string) {
	csrfMiddlewareCfg.tokenFieldName = tokenFieldName
}

type CsrfMiddleware struct{}

func Csrf() *CsrfMiddleware {
	return &CsrfMiddleware{}
}

func (a *CsrfMiddleware) Next(ctx *Ctx) (Response, error) {
	csrf := ctx.CsrfToken()

	inputCsrf := ""

	if ctx.Form().Values.Has(csrfMiddlewareCfg.tokenFieldName) {
		inputCsrf = ctx.Form().Values.Get(csrfMiddlewareCfg.tokenFieldName)
	}

	if ctx.Request().URL.Query().Has(csrfMiddlewareCfg.tokenFieldName) {
		inputCsrf = ctx.Request().URL.Query().Get(csrfMiddlewareCfg.tokenFieldName)
	}

	inputCsrf = strings.TrimSpace(inputCsrf)

	ctx.Session().Set(csrfMiddlewareCfg.tokenFieldName, utils.UUID())

	if inputCsrf == "" || inputCsrf != csrf {
		return ctx.CodeResponse(http.StatusForbidden), nil
	}

	return ctx.Next()
}

func (a *CsrfMiddleware) Priority() uint {
	return 200
}
