package vla

import (
	"net/http"
)

type Handle func(w http.ResponseWriter, r *http.Request, vr Route, p Params)

type Middleware func(next Handle) Handle

type Router interface {
	Group

	http.Handler

	Routes() []Route
}

type Group interface {
	Parent() Group
	Router() Router

	Group(prefix string) Group
	Use(m Middleware) Group

	Handle(method, path string, handle Handle) Route
}

type Route interface {
	Parent() Group
	Router() Router
	Method() string

	Path() string
	Name() string
	SetName(name string) Route
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (p Params) ByName(name string) string {
	for i := range p {
		if p[i].Key == name {
			return p[i].Value
		}
	}
	return ""
}

func DELETE(r Router, path string, h Handle) Route {
	return r.Handle(http.MethodDelete, path, h)
}

func GET(r Router, path string, h Handle) Route {
	return r.Handle(http.MethodGet, path, h)
}

func HEAD(r Router, path string, h Handle) Route {
	return r.Handle(http.MethodHead, path, h)
}

func OPTIONS(r Router, path string, h Handle) Route {
	return r.Handle(http.MethodOptions, path, h)
}

func PATCH(r Router, path string, h Handle) Route {
	return r.Handle(http.MethodPatch, path, h)
}

func POST(r Router, path string, h Handle) Route {
	return r.Handle(http.MethodPost, path, h)
}

func PUT(r Router, path string, h Handle) Route {
	return r.Handle(http.MethodPut, path, h)
}
