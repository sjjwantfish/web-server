package web

import (
	"fmt"
	"reflect"
	"strings"
)

// 支持路由树操作
type router struct {
	// http method -> 路由树根节点
	trees map[string]*node
}

func (r *router) AddRoute(method string, path string, handle HandleFunc) {
	// 可增加 method && handle 的校验
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
	if root.handler != nil {
		panic(fmt.Sprintf("重复注册 %s\n", path))
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

func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	path = strings.Trim(path, "/")
	for _, seg := range strings.Split(path, "/") {
		child, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		root = child
	}
	return root, root.handler != nil
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

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}
	child, ok := n.children[path]
	return child, ok
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
