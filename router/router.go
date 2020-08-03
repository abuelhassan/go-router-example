package router

import (
	"net/http"
)

type Router struct {
	routes   []route
	NotFound http.HandlerFunc
}

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

func New() *Router {
	return &Router{
		routes:   []route{},
		NotFound: http.NotFound,
	}
}

func (rtr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path

	var matching http.HandlerFunc
	for _, route := range rtr.routes {
		if method == route.method && path == route.pattern {
			matching = route.handler
			break
		}
	}
	if matching == nil {
		matching = rtr.NotFound
	}
	matching.ServeHTTP(w, r)
}

func (rtr *Router) Route(method, pattern string, handlerFunc http.HandlerFunc) {
	rtr.routes = append(rtr.routes, route{method, pattern, handlerFunc})
}
