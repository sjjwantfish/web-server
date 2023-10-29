package web

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
	Addr string
	*router
}

var _ Server = &HTTPServer{}

func NewHttpServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Req:  req,
		Resp: writer,
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	// 查找路由，命中逻辑
	node, ok := h.findRoute(ctx.Req.Host, ctx.Req.URL.Path)
	if !ok || node.handler == nil {
		// 404
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("details not found"))
		return
	}
	node.handler(ctx)
}

func (h *HTTPServer) Start() error {
	l, err := net.Listen("tcp", h.Addr)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}

// func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {

// }

type HTTPSServer struct {
	HTTPServer
}

var _ Server = &HTTPSServer{}

type HandleFunc func(ctx *Context)
