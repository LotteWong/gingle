package gingle

import "net/http"

type router struct {
	handlers map[string]HandlerFunc // TODO: 路由映射哈希
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// TODO: 注册路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// TODO: 处理路由
func (r *router) handleRoute(ctx *Context) {
	key := ctx.Method + "-" + ctx.Pattern
	if handler, ok := r.handlers[key]; ok {
		handler(ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 Not Found: %s\n", ctx.Pattern)
	}
}
