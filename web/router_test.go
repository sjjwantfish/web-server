package web

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_router_AddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{method: "GET", path: "/user/home"},
		{method: "GET", path: "/"},
		{method: "POST", path: "/api/user/login"},
		{method: "POST", path: "/api/user/logout"},
		{method: "POST", path: "/api/user"},
		{method: "POST", path: "login"},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
					},
				},
				handler: mockHandler,
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"api": {
						path: "api",
						children: map[string]*node{
							"user": {
								path: "user",
								children: map[string]*node{
									"login": {
										path:    "login",
										handler: mockHandler,
									},
								},
								handler: mockHandler,
							},
						},
					},
					"login": {
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)
}

func Test_router_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{method: "GET", path: "/user/home"},
		{method: "GET", path: "/"},
		{method: "POST", path: "/api/user/login"},
		{method: "POST", path: "/api/user/logout"},
		{method: "POST", path: "/api/user"},
		{method: "POST", path: "login"},
		{method: "DELETE", path: "/"},
	}
	r := newRouter()
	mockHandler := func(ctx *Context) {}
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNode  *node
	}{
		{
			name:      "method not found",
			method:    "OPTIONS",
			path:      "/user/home",
			wantFound: false,
			wantNode: &node{
				path:    "detail",
				handler: mockHandler,
			},
		},
		{
			name:      "user home",
			method:    "GET",
			path:      "/user/home",
			wantFound: true,
			wantNode: &node{
				path:    "home",
				handler: mockHandler,
			},
		},
		{
			name:      "no handler",
			method:    "GET",
			path:      "/user",
			wantFound: false,
			wantNode: &node{
				path:    "user",
				handler: nil,
			},
		},
		{
			name:      "root",
			method:    "DELETE",
			path:      "/",
			wantFound: false,
			wantNode: &node{
				path:    "/",
				handler: nil,
			},
		},
		{
			name:      "path no found",
			method:    "GET",
			path:      "/user/home1",
			wantFound: false,
			wantNode: &node{
				path:    "home1",
				handler: mockHandler,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.wantNode.path, n.path)
			assert.Equal(t, tc.wantNode.children, n.children)

			nHandler := reflect.ValueOf(n.handler).Pointer()
			yHandler := reflect.ValueOf(tc.wantNode.handler).Pointer()
			assert.True(t, nHandler == yHandler)
		})
	}
}
