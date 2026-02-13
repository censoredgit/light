package controller

import "errors"

var authMiddlewareCfg struct {
	loginRouteUri           string
	loginRouteName          string
	backAfterAuthQueryParam string
}

type AuthMiddleware struct{}

func setupAuthMiddleware(loginRouteUri, loginRouteName, backAfterAuthQueryParam string) {
	authMiddlewareCfg.loginRouteUri = loginRouteUri
	authMiddlewareCfg.loginRouteName = loginRouteName
	authMiddlewareCfg.backAfterAuthQueryParam = backAfterAuthQueryParam
}

func Auth() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (a *AuthMiddleware) Next(ctx *Ctx) (Response, error) {
	if authMiddlewareCfg.loginRouteUri == "" {
		return nil, errors.New("login uri required")
	}

	if ctx.IsGuest() {
		ref := ctx.Request().Referer()
		if ref == "" && ctx.Request().RequestURI != authMiddlewareCfg.loginRouteUri {
			ref = ctx.Request().RequestURI
		}
		if ref != "" {
			ctx.Session().Set(authMiddlewareCfg.backAfterAuthQueryParam, ref)
		}
		return ctx.RedirectResponse(authMiddlewareCfg.loginRouteUri), nil
	}

	return ctx.Next()
}

func (a *AuthMiddleware) Priority() uint {
	return 1
}
