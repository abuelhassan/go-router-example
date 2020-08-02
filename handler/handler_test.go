package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Run("alive", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		HealthCheck(w, r)

		actualResponse, err := ioutil.ReadAll(w.Result().Body)
		assert.Nil(t, err)
		assert.Equal(t, "{\"alive\":true}", string(actualResponse))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestNotFound(t *testing.T) {
	t.Run("route not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		NotFound(w, r)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}
