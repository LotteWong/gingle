package gingle

import (
	"html/template"
	"net/http"
	"strings"
)

// TODO: Handler 分离路由部件
type Engine struct {
	*RouterGroup
	groups    []*RouterGroup
	router    *router
	templates *template.Template
	funcMap   template.FuncMap
}

func New() *Engine {
	// TODO: 你中有我我中有你
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// TODO: ListenAndServe 监听
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// TODO: ServeHTTP 服务
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	middlewares := make([]HandlerFunc, 0)
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	ctx := newContext(rw, req) // TODO: 关键转换
	ctx.middlewares = middlewares

	engine.router.handle(ctx) // TODO: 关键转换
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.templates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
