package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRouter_ServeHTTP(t *testing.T) {
	originalNotFound := defaultNotFoundHandler
	defer func() {
		defaultNotFoundHandler = originalNotFound
	}()
	type globals struct {
		defaultNotFound http.HandlerFunc
	}
	type fields struct {
		NotFoundHandler http.Handler
		routes          []route
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
			name: "match pattern and method",
			globals: globals{
				defaultNotFound: nil,
			},
			fields: fields{
				NotFoundHandler: nil,
				routes: []route{
					{
						method:  http.MethodGet,
						pattern: "/mismatch",
						handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(200)
							_, _ = w.Write([]byte("mismatch"))
						}),
					},
					{
						method:  http.MethodPost,
						pattern: "/match",
						handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(200)
							_, _ = w.Write([]byte("post match"))
						}),
					},
					{
						method:  http.MethodGet,
						pattern: "/match",
						handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(200)
							_, _ = w.Write([]byte("get match"))
						}),
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/match", nil),
			},
			wantStatus: 200,
			wantBody:   "get match",
		},
		{
			name: "handler not found",
			globals: globals{
				defaultNotFound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
					_, _ = w.Write([]byte("default handler not found"))
				}),
			},
			fields: fields{
				NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
					_, _ = w.Write([]byte("handler not found"))
				}),
				routes: []route{
					{
						method:  http.MethodGet,
						pattern: "/match",
						handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(200)
							_, _ = w.Write([]byte("get mismatch"))
						}),
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/mismatch", nil),
			},
			wantStatus: 404,
			wantBody:   "handler not found",
		},
		{
			name: "default handler not found",
			globals: globals{
				defaultNotFound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
					_, _ = w.Write([]byte("default handler not found"))
				}),
			},
			fields: fields{
				NotFoundHandler: nil,
				routes: []route{
					{
						method:  http.MethodGet,
						pattern: "/match",
						handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(200)
							_, _ = w.Write([]byte("get mismatch"))
						}),
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/mismatch", nil),
			},
			wantStatus: 404,
			wantBody:   "default handler not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultNotFoundHandler = tt.globals.defaultNotFound
			rtr := Router{
				NotFoundHandler: tt.fields.NotFoundHandler,
				routes:          tt.fields.routes,
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

func TestRouter_matchRequest(t *testing.T) {
	type fields struct {
		routes []route
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    route
		wantErr error
	}{
		{
			name: "match method and pattern",
			fields: fields{
				routes: []route{
					{
						method:  http.MethodGet,
						pattern: "/mismatch",
					},
					{
						method:  http.MethodPost,
						pattern: "/match",
					},
					{
						method:  http.MethodGet,
						pattern: "/match",
					},
				},
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/match", nil),
			},
			want: route{
				method:  http.MethodGet,
				pattern: "/match",
			},
			wantErr: nil,
		},
		{
			name: "no match",
			fields: fields{
				routes: []route{
					{
						method:  http.MethodGet,
						pattern: "/mismatch",
					},
					{
						method:  http.MethodPost,
						pattern: "/match",
					},
				},
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/match", nil),
			},
			want:    route{},
			wantErr: errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtr := &Router{
				routes: tt.fields.routes,
			}
			got, err := rtr.matchRequest(tt.args.r)
			if err != tt.wantErr {
				t.Errorf("matchRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
