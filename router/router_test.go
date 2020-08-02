package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	setup := func(method, target string) (*httptest.ResponseRecorder, *http.Request, *bool, http.HandlerFunc) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, target, nil)

		isInvoked := false
		h := func(http.ResponseWriter, *http.Request) {
			isInvoked = true
		}

		return w, r, &isInvoked, h
	}
	t.Run("match", func(t *testing.T) {
		w, r, isInvoked, h := setup(http.MethodGet, "/match")

		routes := []Route{
			{
				Method:  http.MethodGet,
				Pattern: "/match",
				Handler: h,
			},
		}

		rtr := New(routes, nil)
		rtr.ServeHTTP(w, r)

		assert.True(t, *isInvoked)
	})
	t.Run("no match", func(t *testing.T) {
		w, r, isInvoked, h := setup(http.MethodGet, "/match")

		routes := []Route{
			{
				Method:  http.MethodGet,
				Pattern: "/health",
				Handler: nil,
			},
		}

		rtr := New(routes, h)
		rtr.ServeHTTP(w, r)

		assert.True(t, *isInvoked)
	})
}
