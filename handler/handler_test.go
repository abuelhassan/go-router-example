package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceHandlers(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		handlerName string
		handler     http.HandlerFunc
		args        args
		wantStatus  int
		wantBody    string
	}{
		{
			name:        "health check",
			handlerName: "HealthCheck()",
			handler:     HealthCheck,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/health", nil),
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"alive":true}`,
		},
		{
			name:        "not found",
			handlerName: "NotFound()",
			handler:     NotFound,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "",
		},
		{
			name:        "method not allowed",
			handlerName: "MethodNotAllowed()",
			handler:     MethodNotAllowed,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler(tt.args.w, tt.args.r)

			res := tt.args.w.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {
					t.Fatalf("%s couldn't close body", tt.handlerName)
				}
			}()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("%s statusCode = %d, wantStatusCode %d", tt.handlerName, res.StatusCode, tt.wantStatus)
			}

			gotBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("%s couldn't read body, wantBody %s", tt.handlerName, gotBody)
			}
			if string(gotBody) != tt.wantBody {
				t.Errorf("%s body = %s, wantBody %s", tt.handlerName, gotBody, tt.wantBody)
			}
		})
	}
}
