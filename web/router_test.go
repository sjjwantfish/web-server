package web

import (
	"fmt"
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
		{method: "POST", path: "/api/user/login"},
		{method: "POST", path: "/api/user/logout"},
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
							},
						},
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)
}

func (n *router) equal(y *router) (string, bool) {
	for k, v := range n.trees {
		if dst, ok := y.trees[k]; !ok {
			return fmt.Sprintf("trees: %s not found", k), false
		} else {
			msg, equal := v.equal(dst)
			return msg, equal
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {

	if n.path != y.path {
		return fmt.Sprint("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprint("子节点数量不相等"), false
	}

	nHandler := reflect.ValueOf(n.handler).Pointer()
	yHandler := reflect.ValueOf(y.handler).Pointer()
	fmt.Printf("%v %v\n", nHandler, yHandler)
	if nHandler != yHandler {
		return fmt.Sprintf("%s handler 不相等", n.path), false
	}

	for path, child := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}
		msg, ok := dst.equal(child)
		if !ok {
			return msg, false
		}
	}
	return "", true
}
