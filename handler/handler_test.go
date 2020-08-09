package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "route not found",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/health", nil),
			},
			wantStatus: 200,
			wantBody:   `{"alive":true}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HealthCheck(tt.args.w, tt.args.r)

			res := tt.args.w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("HealthCheck() statusCode = %d, wantStatusCode %d", res.StatusCode, tt.wantStatus)
			}

			gotBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("HealthCheck() couldn't read body, wantBody %s", gotBody)
			}
			if string(gotBody) != tt.wantBody {
				t.Errorf("HealthCheck() body = %s, wantBody %s", gotBody, tt.wantBody)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "route not found",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatus: 404,
			wantBody:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NotFound(tt.args.w, tt.args.r)

			res := tt.args.w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("NotFound() statusCode = %d, wantStatusCode %d", res.StatusCode, tt.wantStatus)
			}

			gotBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("NotFound() couldn't read body, wantBody %s", gotBody)
			}
			if string(gotBody) != tt.wantBody {
				t.Errorf("NotFound() body = %s, wantBody %s", gotBody, tt.wantBody)
			}
		})
	}
}
