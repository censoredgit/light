package controller

import "net/url"

type Response interface {
	WithInputs()
	ExceptInputs(fields ...string)
	OnlyInputs(fields ...string)
	Process(ctx *Ctx)
	With(func(response ResponseExtendData)) Response
	HasData() bool
	Data() any
}

type ResponseExtendData interface {
	SetData(data any)
	AddMessage(key string, data string)
	AddError(key string, data string)
	SetErrors(map[string]string)
	SetOldInput(val url.Values)
	SetInput(val url.Values)
}
