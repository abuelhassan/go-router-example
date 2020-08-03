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

		rtr := New()
		rtr.Route(http.MethodGet, "/match", h)
		rtr.ServeHTTP(w, r)

		assert.True(t, *isInvoked)
	})
	t.Run("no match", func(t *testing.T) {
		w, r, notFoundInvoked, h := setup(http.MethodGet, "/match")
		handlerInvoked := false

		rtr := New()
		rtr.Route(http.MethodGet, "/health", func(writer http.ResponseWriter, request *http.Request) {
			handlerInvoked = true
		})
		rtr.NotFound = h
		rtr.ServeHTTP(w, r)

		assert.False(t, handlerInvoked)
		assert.True(t, *notFoundInvoked)
	})
	t.Run("default no match", func(t *testing.T) {
		w, r, isInvoked, h := setup(http.MethodGet, "/match")

		rtr := New()
		rtr.Route(http.MethodGet, "/health", h)
		rtr.ServeHTTP(w, r)

		assert.False(t, *isInvoked)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}
