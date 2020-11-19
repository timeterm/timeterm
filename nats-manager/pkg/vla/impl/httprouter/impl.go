package vlahttprouter

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/timeterm/timeterm/nats-manager/pkg/vla"
)

type Router struct {
	router      *httprouter.Router
	middlewares []vla.Middleware

	routes []vla.Route
}

func New() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

func (r *Router) newRoute(method, path string) *Route {
	return newRoute(r, r, r.middlewares, method, path)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, hr *http.Request) {
	r.router.ServeHTTP(w, hr)
}

func (r *Router) Routes() []vla.Route {
	return r.routes
}

func (r *Router) Parent() vla.Group {
	return r
}

func (r *Router) Router() vla.Router {
	return r
}

func (r *Router) Group(prefix string) vla.Group {
	return &Group{
		router:      r,
		parent:      r,
		middlewares: r.middlewares,
		prefix:      prefix,
	}
}

func (r *Router) Use(m vla.Middleware) vla.Group {
	r.middlewares = append(r.middlewares, m)
	return r
}

func (r *Router) DELETE(path string, handle vla.Handle) vla.Route {
	route := r.newRoute(http.MethodDelete, path)
	r.router.DELETE(path, newRouteHandle(route, handle))
	return route
}

func (r *Router) GET(path string, handle vla.Handle) vla.Route {
	route := r.newRoute(http.MethodGet, path)
	r.router.GET(path, newRouteHandle(route, handle))
	return route
}

func (r *Router) HEAD(path string, handle vla.Handle) vla.Route {
	route := r.newRoute(http.MethodHead, path)
	r.router.HEAD(path, newRouteHandle(route, handle))
	return route
}

func (r *Router) OPTIONS(path string, handle vla.Handle) vla.Route {
	route := r.newRoute(http.MethodOptions, path)
	r.router.OPTIONS(path, newRouteHandle(route, handle))
	return route
}

func (r *Router) PATCH(path string, handle vla.Handle) vla.Route {
	route := r.newRoute(http.MethodPatch, path)
	r.router.PATCH(path, newRouteHandle(route, handle))
	return route
}

func (r *Router) POST(path string, handle vla.Handle) vla.Route {
	route := r.newRoute(http.MethodPost, path)
	r.router.POST(path, newRouteHandle(route, handle))
	return route
}

func (r *Router) PUT(path string, handle vla.Handle) vla.Route {
	route := r.newRoute(http.MethodPut, path)
	r.router.PUT(path, newRouteHandle(route, handle))
	return route
}

type Group struct {
	router      *Router
	parent      vla.Group
	middlewares []vla.Middleware

	prefix string
}

func (g *Group) newRoute(method, path string) *Route {
	return newRoute(g.router, g, g.middlewares, method, joinPaths(g.prefix, path))
}

func (g *Group) Parent() vla.Group {
	return g.parent
}

func (g *Group) Router() vla.Router {
	return g.router
}

func (g *Group) Group(prefix string) vla.Group {
	return &Group{
		router:      g.router,
		parent:      g,
		middlewares: g.middlewares,
		prefix:      joinPaths(g.prefix, prefix),
	}
}

func (g *Group) Use(m vla.Middleware) vla.Group {
	g.middlewares = append(g.middlewares, m)
	return g
}

func (g *Group) DELETE(path string, handle vla.Handle) vla.Route {
	route := g.newRoute(http.MethodDelete, path)
	g.router.router.DELETE(joinPaths(g.prefix, path), newRouteHandle(route, handle))
	return route
}

func (g *Group) GET(path string, handle vla.Handle) vla.Route {
	route := g.newRoute(http.MethodGet, path)
	g.router.router.GET(joinPaths(g.prefix, path), newRouteHandle(route, handle))
	return route
}

func (g *Group) HEAD(path string, handle vla.Handle) vla.Route {
	route := g.newRoute(http.MethodHead, path)
	g.router.router.HEAD(joinPaths(g.prefix, path), newRouteHandle(route, handle))
	return route
}

func (g *Group) OPTIONS(path string, handle vla.Handle) vla.Route {
	route := g.newRoute(http.MethodOptions, path)
	g.router.router.OPTIONS(joinPaths(g.prefix, path), newRouteHandle(route, handle))
	return route
}

func (g *Group) PATCH(path string, handle vla.Handle) vla.Route {
	route := g.newRoute(http.MethodPatch, path)
	g.router.router.PATCH(joinPaths(g.prefix, path), newRouteHandle(route, handle))
	return route
}

func (g *Group) POST(path string, handle vla.Handle) vla.Route {
	route := g.newRoute(http.MethodPost, path)
	g.router.router.POST(joinPaths(g.prefix, path), newRouteHandle(route, handle))
	return route
}

func (g *Group) PUT(path string, handle vla.Handle) vla.Route {
	route := g.newRoute(http.MethodPut, path)
	g.router.router.PUT(joinPaths(g.prefix, path), newRouteHandle(route, handle))
	return route
}

type Route struct {
	router      *Router
	parent      vla.Group
	middlewares []vla.Middleware

	method string
	path   string
	name   string
}

func newRoute(r *Router, p vla.Group, m []vla.Middleware, method, path string) *Route {
	route := &Route{
		router:      r,
		parent:      p,
		middlewares: m,
		method:      method,
		path:        path,
	}
	r.routes = append(r.routes, route)

	return route
}

func (r *Route) Parent() vla.Group {
	return r.parent
}

func (r *Route) Router() vla.Router {
	return r.router
}

func (r *Route) Method() string {
	return r.method
}

func (r *Route) Path() string {
	return r.path
}

func (r *Route) Name() string {
	return r.name
}

func (r *Route) SetName(name string) vla.Route {
	r.name = name
	return r
}

func convertParams(params httprouter.Params) vla.Params {
	vp := make(vla.Params, len(params))
	for i := range params {
		vp[i] = vla.Param(params[i])
	}
	return vp
}

func joinPaths(a, b string) string {
	pat := strings.Join([]string{a, b}, "/")
	for strings.Contains(pat, "//") {
		pat = strings.ReplaceAll(pat, "//", "/")
	}
	return pat
}

func withParamsInContext(next vla.Handle) vla.Handle {
	return func(w http.ResponseWriter, r *http.Request, vr vla.Route, p vla.Params) {
		ctx := vla.ContextWithParams(r.Context(), p)
		next(w, r.WithContext(ctx), vr, p)
	}
}

func withRouteInContext(next vla.Handle) vla.Handle {
	return func(w http.ResponseWriter, r *http.Request, vr vla.Route, p vla.Params) {
		ctx := vla.ContextWithRoute(r.Context(), vr)
		next(w, r.WithContext(ctx), vr, p)
	}
}

func newRouteHandle(vr *Route, next vla.Handle) httprouter.Handle {
	middlewares := append([]vla.Middleware{
		withParamsInContext,
		withRouteInContext,
	}, vr.middlewares...)

	hdl := useMiddlewares(middlewares, next)
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		hdl(w, r, vr, convertParams(p))
	}
}

func useMiddlewares(m []vla.Middleware, next vla.Handle) vla.Handle {
	hdl := next
	for _, middleware := range m {
		hdl = middleware(hdl)
	}
	return hdl
}
