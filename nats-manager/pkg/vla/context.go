package vla

import "context"

type paramsKeyType struct{}

var paramsKey paramsKeyType

func ParamsFromContext(ctx context.Context) Params {
	params, ok := ctx.Value(paramsKey).(Params)
	if !ok {
		return nil
	}
	return params
}

func ContextWithParams(ctx context.Context, p Params) context.Context {
	return context.WithValue(ctx, paramsKey, p)
}

type routeKeyType struct{}

var routeKey routeKeyType

func RouteFromContext(ctx context.Context) Route {
	route, ok := ctx.Value(routeKey).(Route)
	if !ok {
		return nil
	}
	return route
}

func ContextWithRoute(ctx context.Context, r Route) context.Context {
	return context.WithValue(ctx, routeKey, r)
}
