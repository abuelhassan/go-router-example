package router

import (
	"errors"
	"net/http"
)

var (
	defaultNotFoundHandler = http.NotFound

	handleNotFound      = notFound
	matchRequestHandler = matchRequest

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
	matching, err := matchRequestHandler(rtr.routes, r)
	if err != nil {
		switch err {
		case errNotFound:
			handleNotFound(rtr.NotFoundHandler, w, r)
		}
		return
	}
	matching.ServeHTTP(w, r)
}

func (rtr *Router) Route(method, pattern string, handler http.Handler) {
	rtr.routes = append(rtr.routes, route{method, pattern, handler})
}

func matchRequest(routes []route, r *http.Request) (http.Handler, error) {
	method := r.Method
	path := r.URL.Path

	for _, route := range routes {
		if method == route.method && path == route.pattern {
			return route.handler, nil
		}
	}

	return nil, errNotFound
}

func notFound(h http.Handler, w http.ResponseWriter, r *http.Request) {
	if h != nil {
		h.ServeHTTP(w, r)
		return
	}
	defaultNotFoundHandler(w, r)
}
