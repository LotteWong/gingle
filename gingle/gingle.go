package gingle

import (
	"net/http"
)

// TODO: HandlerFunc 业务处理函数
type HandlerFunc func(*Context)

// TODO: Handler 分离路由部件
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

// TODO: ListenAndServe 监听
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// TODO: ServeHTTP 服务
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := newContext(rw, req) // TODO: 关键转换
	engine.router.handle(ctx)  // TODO: 关键转换
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.addRoute("POST", pattern, handler)
}

func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.router.addRoute("PUT", pattern, handler)
}

func (engine *Engine) DELETE(pattern string, handler HandlerFunc) {
	engine.router.addRoute("DELETE", pattern, handler)
}
