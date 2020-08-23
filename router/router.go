package router

import (
	"errors"
	"net/http"

	"github.com/abuelhassan/go-router-example/trie"
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

type (
	Router struct {
		NotFoundHandler         http.Handler
		MethodNotAllowedHandler http.Handler
		routes                  trie.Trier
	}

	// route is a map from http.Method to http.Handler.
	route map[string]http.Handler
)

// New returns an instance of Router
func New() Router {
	return Router{
		routes: trie.NewPathTrie(),
	}
}

func (rtr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, err := matchRouteFunc(rtr.routes, r)
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
	h.ServeHTTP(w, r)
}

// Route adds a new route, or overrides it if it already exists.
func (rtr *Router) Route(method, pattern string, handler http.Handler) {
	val := route{}
	r := rtr.routes.Get(pattern)
	if r != nil {
		val = r.(route)
	}
	val[method] = handler
	rtr.routes.Put(pattern, val)
}

func matchRoute(trie trie.Trier, r *http.Request) (http.Handler, error) {
	method := r.Method
	path := r.URL.Path

	v := trie.Get(path)
	if v == nil {
		return nil, errNotFound
	}

	mp, ok := v.(route)
	if !ok || len(mp) == 0 {
		// TODO: log error. maybe use wrapping!
		return nil, errNotFound
	}

	if mp[method] == nil {
		return nil, errMethodNotAllowed
	}

	return mp[method], nil
}
