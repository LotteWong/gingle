package gingle

import (
	"fmt"
	"net/http"
)

// TODO: HandlerFunc 业务处理函数
type HandlerFunc func(http.ResponseWriter, *http.Request)

// TODO: Handler 路由映射哈希
type Engine struct {
	router map[string]HandlerFunc
}

// TODO: 新建
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

// TODO: ListenAndServe 监听
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// TODO: ServeHTTP 服务
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(rw, req)
	} else {
		fmt.Fprintf(rw, "404 NOT FOUND: %s\n", req.URL)
	}
}

// TODO: 注册
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.addRoute("PUT", pattern, handler)
}

func (engine *Engine) DELETE(pattern string, handler HandlerFunc) {
	engine.addRoute("DELETE", pattern, handler)
}
