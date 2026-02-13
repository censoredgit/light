package controller

import (
	"fmt"
	"path"
	"strings"
)

type MountInfo struct {
	action        *Action
	middlewares   []Middleware
	routeName     string
	routeUri      string
	method        string
	hasParameters bool
}

func (m *MountInfo) Name(name string) {
	m.routeName = name
}

type Mount struct {
	info   map[string]*MountInfo
	di     *Container
	belong *Group
	groups []*Group
}

func newMount(c *Container) *Mount {
	return &Mount{info: make(map[string]*MountInfo), di: c}
}

func (m *Mount) Get(uri string, fa any, middlewares ...Middleware) *MountInfo {
	return m.mount("GET", uri, m.di.inject(fa).(*Action), middlewares...)
}

func (m *Mount) Post(uri string, fa any, middlewares ...Middleware) *MountInfo {
	return m.mount("POST", uri, m.di.inject(fa).(*Action), middlewares...)
}

func (m *Mount) Delete(uri string, fa any, middlewares ...Middleware) *MountInfo {
	return m.mount("DELETE", uri, m.di.inject(fa).(*Action), middlewares...)
}

func (m *Mount) Put(uri string, fa any, middlewares ...Middleware) *MountInfo {
	return m.mount("PUT", uri, m.di.inject(fa).(*Action), middlewares...)
}

func (m *Mount) Patch(uri string, fa any, middlewares ...Middleware) *MountInfo {
	return m.mount("PATCH", uri, m.di.inject(fa).(*Action), middlewares...)
}

func (m *Mount) Group() *Group {
	g := &Group{
		parent: m.belong,
	}
	g.m = newMount(m.di)
	g.m.belong = g
	m.groups = append(m.groups, g)
	return g
}

func (m *Mount) mount(method, uri string, action *Action, middlewares ...Middleware) *MountInfo {
	uriWithMethod := fmt.Sprintf("%s %s", method, uri)
	if m.info[uriWithMethod] == nil {
		m.info[uriWithMethod] = &MountInfo{action: action, method: method, hasParameters: strings.Contains(uri, "{")}

		action.name = &(m.info[uriWithMethod]).routeName
		action.uri = &(m.info[uriWithMethod]).routeUri
		action.WithMiddleware(middlewares...)

	} else {
		panic(fmt.Sprintf("action already defined for uri %s %s", method, path.Join(m.belong.composePrefix(), uri)))
	}

	return m.info[uriWithMethod]
}
