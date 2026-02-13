package controller

import "strings"

type Group struct {
	parent      *Group
	m           *Mount
	prefix      string
	parameters  bool
	middlewares []Middleware
}

func (g *Group) hasParameters() bool {
	hasParameters := g.parameters
	if g.parent != nil && g.parent != g {
		hasParameters = hasParameters || g.parent.hasParameters()
	}

	return hasParameters
}

func (g *Group) composePrefix() string {
	prefix := g.prefix
	if g.parent != nil && g.parent != g {
		prefix = g.parent.composePrefix() + prefix
	}
	return prefix
}

func (g *Group) composeMiddleware() []Middleware {
	middlewares := make([]Middleware, 0)
	if g.parent != nil && g.parent != g {
		middlewares = append(middlewares, g.parent.composeMiddleware()...)
	}
	middlewares = append(middlewares, g.middlewares...)
	return middlewares
}

func (g *Group) Prefix(prefix string) *Group {
	g.prefix = prefix
	g.parameters = strings.Contains(prefix, "{")

	return g
}

func (g *Group) Middlewares(middlewares ...Middleware) *Group {
	g.middlewares = append(g.middlewares, middlewares...)

	return g
}

func (g *Group) Mount(f func(m *Mount)) {
	f(g.m)
}
