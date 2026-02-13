package controller

import (
	"net/http"
)

type PermissionMiddleware struct {
	permissions []string
}

func Permission(permissions ...string) *PermissionMiddleware {
	return &PermissionMiddleware{permissions: permissions}
}

func (a *PermissionMiddleware) Next(ctx *Ctx) (Response, error) {
	hasRole := false
	for _, permission := range a.permissions {
		if ctx.Can(permission) {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return ctx.CodeResponse(http.StatusUnauthorized), nil
	}

	return ctx.Next()
}

func (a *PermissionMiddleware) Priority() uint {
	return 300
}
