package vla

import "context"

type paramsKeyType struct{}

func ParamsFromContext(ctx context.Context) Params {
	params, ok := ctx.Value(paramsKeyType{}).(Params)
	if !ok {
		return nil
	}
	return params
}

func ContextWithParams(ctx context.Context, p Params) context.Context {
	return context.WithValue(ctx, paramsKeyType{}, p)
}

type routeKeyType struct{}

func RouteFromContext(ctx context.Context) Route {
	route, ok := ctx.Value(routeKeyType{}).(Route)
	if !ok {
		return nil
	}
	return route
}

func ContextWithRoute(ctx context.Context, r Route) context.Context {
	return context.WithValue(ctx, routeKeyType{}, r)
}
