package gingle

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // TODO: 路由映射节点
	handlers map[string]HandlerFunc // TODO: 路由映射函数
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// TODO: 自动补全问题
func convertPatternToParts(pattern string) []string {
	items := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// TODO: 注册方法问题
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern

	parts := convertPatternToParts(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)

	r.handlers[key] = handler
}

// TODO: 处理方法问题
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := convertPatternToParts(pattern)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		foundParts := convertPatternToParts(n.pattern)
		for index, part := range foundParts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) >= 2 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}

	nodes := make([]*node, 0)
	root.traverse(&nodes) // TODO: 这里应该是方法路由的集合
	return nodes
}

func (r *router) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Pattern)

	if n != nil {
		ctx.Params = params
		key := ctx.Method + "-" + n.pattern // TODO: 注意是n.pattern不是ctx.Pattern
		r.handlers[key](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Pattern)
	}
}
