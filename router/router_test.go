package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/abuelhassan/go-router-example/trie"
)

func TestRouter_ServeHTTP(t *testing.T) {
	originalNotFound := defaultNotFoundHandler
	originalMethodNotAllowed := defaultMethodNotAllowedHandler
	originalMatchRouteFunc := matchRouteFunc
	defer func() {
		defaultNotFoundHandler = originalNotFound
		defaultMethodNotAllowedHandler = originalMethodNotAllowed
		matchRouteFunc = originalMatchRouteFunc
	}()
	defaultNotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("default not found"))
	}
	defaultMethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("default method not allowed"))
	}
	type globals struct {
		matchRouteFunc func(trie.Trier, *http.Request) (http.Handler, error)
	}
	type fields struct {
		NotFoundHandler         http.Handler
		MethodNotAllowedHandler http.Handler
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		globals    globals
		fields     fields
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "match route",
			globals: globals{
				matchRouteFunc: func(_ trie.Trier, request *http.Request) (http.Handler, error) {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte("match found"))
					}), nil
				},
			},
			fields: fields{
				NotFoundHandler:         nil,
				MethodNotAllowedHandler: nil,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/match", nil),
			},
			wantStatus: http.StatusOK,
			wantBody:   "match found",
		},
		{
			name: "not found",
			globals: globals{
				matchRouteFunc: func(_ trie.Trier, request *http.Request) (http.Handler, error) {
					return nil, fmt.Errorf("wrapped error: %w", errNotFound)
				},
			},
			fields: fields{
				NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write([]byte("not found"))
				}),
				MethodNotAllowedHandler: nil,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "not found",
		},
		{
			name: "default not found",
			globals: globals{
				matchRouteFunc: func(_ trie.Trier, request *http.Request) (http.Handler, error) {
					return nil, fmt.Errorf("wrapped error: %w", errNotFound)
				},
			},
			fields: fields{
				NotFoundHandler:         nil,
				MethodNotAllowedHandler: nil,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "default not found",
		},
		{
			name: "method not allowed",
			globals: globals{
				matchRouteFunc: func(_ trie.Trier, request *http.Request) (http.Handler, error) {
					return nil, fmt.Errorf("wrapped error: %w", errMethodNotAllowed)
				},
			},
			fields: fields{
				NotFoundHandler: nil,
				MethodNotAllowedHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusMethodNotAllowed)
					_, _ = w.Write([]byte("method not allowed"))
				}),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name: "default method not allowed",
			globals: globals{
				matchRouteFunc: func(_ trie.Trier, request *http.Request) (http.Handler, error) {
					return nil, fmt.Errorf("wrapped error: %w", errMethodNotAllowed)
				},
			},
			fields: fields{
				NotFoundHandler:         nil,
				MethodNotAllowedHandler: nil,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "default method not allowed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchRouteFunc = tt.globals.matchRouteFunc
			rtr := Router{
				NotFoundHandler:         tt.fields.NotFoundHandler,
				MethodNotAllowedHandler: tt.fields.MethodNotAllowedHandler,
			}

			rtr.ServeHTTP(tt.args.w, tt.args.r)
			got := tt.args.w.Result()
			defer func() {
				err := got.Body.Close()
				if err != nil {
					t.Fatal("ServeHTTP() couldn't close body")
				}
			}()

			if got.StatusCode != tt.wantStatus {
				t.Errorf("ServeHTTP() StatusCode = %d, want %d", got.StatusCode, tt.wantStatus)
			}

			gotBody, err := ioutil.ReadAll(got.Body)
			if err != nil {
				t.Fatalf("ServeHTTP() couldn't read body, wantBody %s", gotBody)
			}
			if string(gotBody) != tt.wantBody {
				t.Errorf("HealthCheck() body = %s, wantBody %s", gotBody, tt.wantBody)
			}
		})
	}
}

func Test_matchRoute(t *testing.T) {
	match := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	type args struct {
		trie trie.Trier
		r    *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    http.Handler
		wantErr error
	}{
		{
			name: "match",
			args: args{
				trie: mockTrier{get: route{http.MethodGet: match}},
				r:    httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want:    match,
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				trie: mockTrier{get: nil},
				r:    httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want:    nil,
			wantErr: fmt.Errorf("%w - conversion error", errNotFound),
		},
		{
			name: "method not allowed",
			args: args{
				trie: mockTrier{get: route{http.MethodGet: match}},
				r:    httptest.NewRequest(http.MethodPost, "/", nil),
			},
			want:    nil,
			wantErr: errMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchRoute(tt.args.trie, tt.args.r)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("matchRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != (tt.want == nil) {
				t.Errorf("matchRoute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
