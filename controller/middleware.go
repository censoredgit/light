package controller

type Middleware interface {
	Next(ctx *Ctx) (Response, error)
	Priority() uint
}
