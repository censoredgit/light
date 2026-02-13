package controller

import (
	"testing"
)

func TestControllerDeepUriPathSuccess(t *testing.T) {
	ctr := newController()
	ctr.Group().Prefix("/a").Mount(func(m *Mount) {
		m.Group().Prefix("/b").Mount(func(m *Mount) {
			m.Get("/c", NewAction(func(ctx *Ctx) (Response, error) {
				return nil, nil
			})).Name("test")
		})
	})

	err := ctr.composeRouters()
	if err != nil {
		t.Error(err)
	}
	if namedRouterMap["test"] != "/a/b/c" {
		t.Error("path should be /a/b/c")
	}
}
