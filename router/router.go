package router

import (
	"errors"
	"net/http"
)

var (
	defaultNotFoundHandler         = http.HandlerFunc(http.NotFound)
	defaultMethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	matchRouteFunc = matchRoute

	errNotFound         = errors.New("not found")
	errMethodNotAllowed = errors.New("method not allowed")
)

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
	routes                  []route
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
	matching, err := matchRouteFunc(rtr.routes, r)
	if err != nil {
		switch err {
		case errNotFound:
			h := rtr.NotFoundHandler
			if h == nil {
				h = defaultNotFoundHandler
			}
			h.ServeHTTP(w, r)
		case errMethodNotAllowed:
			h := rtr.MethodNotAllowedHandler
			if h == nil {
				h = defaultMethodNotAllowedHandler
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

func matchRoute(routes []route, r *http.Request) (route, error) {
	method := r.Method
	path := r.URL.Path

	methodNotAllowed := false
	for _, route := range routes {
		if method == route.method && path == route.pattern {
			return route, nil
		}
		if path == route.pattern {
			methodNotAllowed = true
		}
	}
	if methodNotAllowed {
		return route{}, errMethodNotAllowed
	}
	return route{}, errNotFound
}
