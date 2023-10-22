package server

import (
	"net"
	"net/http"
)

type Server interface {
	http.Handler
	Start() error
	// 路由注册
	AddRoute(method string, path string, handlerFunc HandleFunc)
}

type HTTPServer struct {
	addr string
}

var _ Server = &HTTPServer{}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Req:  req,
		Resp: writer,
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	// 查找路由，命中逻辑

}

func (h *HTTPServer) Start() error {
	l, err := net.Listen("tcp", h.addr)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}

func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {

}

type HTTPSServer struct {
	HTTPServer
}

var _ Server = &HTTPSServer{}

type HandleFunc func(ctx *Context)
