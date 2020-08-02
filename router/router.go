package router

import (
	"context"
	"net/http"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

func New(routes []Route, notFound http.HandlerFunc) http.Handler {
	return router{routes, notFound}
}

type router struct {
	routes   []Route
	notFound http.HandlerFunc
}

func (rtr router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path

	ctx := context.TODO()
	for _, route := range rtr.routes {
		if method != route.Method || path != route.Pattern {
			continue
		}
		route.Handler(w, r.WithContext(ctx))
		return
	}
	rtr.notFound(w, r.WithContext(ctx))
}
