package controller

import "context"

type UserProvider interface {
	GetAuthIdentification(ctx context.Context, authId string) (AuthIdentification, error)
	GetRoleSupport(ctx context.Context, authId string) (RoleSupport, error)
}
