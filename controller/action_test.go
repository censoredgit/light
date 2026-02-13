package controller

import "testing"

func TestActionSuccess(t *testing.T) {
	c := newController()
	c.Get("/", func() *Action {
		return NewAction(func(ctx *Ctx) (Response, error) {
			return nil, nil
		})
	})
}
