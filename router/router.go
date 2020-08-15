package router

import (
	"errors"
	"net/http"
)

var (
	defaultNotFoundHandler = http.NotFound

	errNotFound = errors.New("not found")
)

type Router struct {
	NotFoundHandler http.Handler
	routes          []route
}

type route struct {
	method  string
	pattern string
	handler http.Handler
}

func New() Router {
	return Router{
		routes: []route{},
	}
}

func (rtr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matching, err := rtr.matchRequest(r)
	if err != nil {
		switch err {
		case errNotFound:
			h := rtr.NotFoundHandler
			if h == nil {
				h = http.HandlerFunc(defaultNotFoundHandler)
			}
			h.ServeHTTP(w, r)
		}
		return
	}
	matching.handler.ServeHTTP(w, r)
}

func (rtr *Router) Route(method, pattern string, handler http.Handler) {
	rtr.routes = append(rtr.routes, route{method, pattern, handler})
}

func (rtr *Router) matchRequest(r *http.Request) (route, error) {
	method := r.Method
	path := r.URL.Path

	for _, route := range rtr.routes {
		if method == route.method && path == route.pattern {
			return route, nil
		}
	}

	return route{}, errNotFound
}
