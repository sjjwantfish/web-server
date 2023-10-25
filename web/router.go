package web

import (
	"fmt"
	"strings"
)

// 支持路由树操作
type router struct {
	// http method -> 路由树根节点
	trees map[string]*node
}

func (r *router) AddRoute(method string, path string, handle HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		fmt.Printf("no root of %s\n", method)
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		fmt.Printf("seg: %s\n", seg)
		if seg == "" {
			continue
		}
		root = root.childOrCreate(seg)
	}
	root.handler = handle
	// for _, seg := range segs {
	// 	fmt.Printf("seg: %s \n", seg)
	// 	if seg == "" {
	// 		continue
	//
	// children := root.childOrCreate(seg)
	// root = children
	// }
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

type node struct {
	path string
	// 子 path 到子节点的映射
	children map[string]*node
	// 用户注册的业务逻辑
	handler HandleFunc
}
