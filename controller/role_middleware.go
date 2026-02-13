package controller

import (
	"net/http"
)

type RoleMiddleware struct {
	roles []string
}

func Role(roles ...string) *RoleMiddleware {
	return &RoleMiddleware{roles: roles}
}

func (a *RoleMiddleware) Next(ctx *Ctx) (Response, error) {
	hasRole := false
	for _, role := range a.roles {
		if ctx.Role(role) {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return ctx.CodeResponse(http.StatusUnauthorized), nil
	}

	return ctx.Next()
}

func (a *RoleMiddleware) Priority() uint {
	return 400
}
