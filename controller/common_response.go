package controller

import (
	"net/url"
	"slices"
)

type CommonResponse struct {
	flashStorage ContextFlashStorage
	code         int
	data         any
	withInputs   bool
	exceptInputs []string
	onlyInputs   []string
}

func (c *CommonResponse) SetErrors(m map[string]string) {
	c.flashStorage.Errors().SetErrors(m)
}

func (c *CommonResponse) SetOldInput(val url.Values) {
	c.flashStorage.Inputs().SetOld(val)
}

func (c *CommonResponse) SetInput(val url.Values) {
	c.flashStorage.Inputs().data = val
}

func (c *CommonResponse) AddError(key string, data string) {
	c.flashStorage.Errors().Set(key, data)
}

func (c *CommonResponse) AddMessage(key string, data string) {
	c.flashStorage.Inputs().Set(key, data)
}

func (c *CommonResponse) Process(ctx *Ctx) {
	c.process(ctx)
}

func (c *CommonResponse) WithInputs() {
	c.withInputs = true
}

func (c *CommonResponse) ExceptInputs(fields ...string) {
	c.exceptInputs = fields
}

func (c *CommonResponse) OnlyInputs(fields ...string) {
	c.onlyInputs = fields
}

func (c *CommonResponse) Code() int {
	return c.code
}

func (c *CommonResponse) SetCode(code int) {
	c.code = code
}

func (c *CommonResponse) SetData(data any) {
	c.data = data
}

func (c *CommonResponse) HasData() bool {
	return c.data != nil
}

func (c *CommonResponse) Data() any {
	return c.data
}

func (c *CommonResponse) With(it func(response ResponseExtendData)) Response {
	it(c)
	return c
}

func (c *CommonResponse) process(ctx *Ctx) {
	onlyInputsLen := len(c.onlyInputs)
	exceptInputsLen := len(c.exceptInputs)

	switch {
	case c.withInputs:
		for k, v := range ctx.Form().Values {
			c.flashStorage.Inputs().SetList(k, v)
		}
	case onlyInputsLen > 0:
		for _, field := range c.onlyInputs {
			if ctx.Form().Values.Has(field) {
				c.flashStorage.Inputs().data.Set(field, ctx.Form().Values.Get(field))
			}
		}
	case exceptInputsLen > 0:
		for field, value := range ctx.Form().Values {
			if !slices.Contains(c.exceptInputs, field) {
				c.flashStorage.Inputs().data.Set(field, value[0])
			} else {
				c.flashStorage.Inputs().old.Del(field)
			}
		}
	}
}
