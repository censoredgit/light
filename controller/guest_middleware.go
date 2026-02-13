package controller

type GuestMiddleware struct{}

func Guest() *GuestMiddleware {
	return &GuestMiddleware{}
}

func (a *GuestMiddleware) Next(ctx *Ctx) (Response, error) {
	if ctx.IsAuth() {
		return ctx.RedirectResponse(*rootAction.uri), nil
	}

	return ctx.Next()
}

func (a *GuestMiddleware) Priority() uint {
	return 50
}
