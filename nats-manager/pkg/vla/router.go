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

	DELETE(path string, handle Handle) Route
	GET(path string, handle Handle) Route
	HEAD(path string, handle Handle) Route
	OPTIONS(path string, handle Handle) Route
	PATCH(path string, handle Handle) Route
	POST(path string, handle Handle) Route
	PUT(path string, handle Handle) Route
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
