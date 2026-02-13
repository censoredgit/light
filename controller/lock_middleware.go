package controller

import (
	"github.com/censoredgit/light/locker"
	"net/http"
)

var lockMiddlewareCfg struct {
	locker *locker.Locker
}

func setupLockMiddleware(locker *locker.Locker) {
	lockMiddlewareCfg.locker = locker
}

type LockMiddleware struct{}

func Lock() *LockMiddleware {
	return &LockMiddleware{}
}

func (a *LockMiddleware) Next(ctx *Ctx) (Response, error) {
	if !ctx.IsAuth() {
		return ctx.CodeResponse(http.StatusForbidden), nil
	}

	lockSign := "LockMiddleware_" + ctx.AuthIdentification()
	ctx.extras["LockMiddleware"] = lockSign

	err := lockMiddlewareCfg.locker.SimpleLock(lockSign)
	if err != nil {
		return nil, err
	}
	defer lockMiddlewareCfg.locker.ReleaseSimpleLock(lockSign)

	return ctx.Next()
}

func (a *LockMiddleware) Priority() uint {
	return 100
}
