package gingle

import (
	"html/template"
	"net/http"
	"strings"
)

// Engine is the core of gingle
type Engine struct {
	// For group control
	*RouterGroup
	groups []*RouterGroup
	// For router mapping
	router *router
	// For template render
	templates *template.Template
	funcMap   template.FuncMap
}

// New returns a instance of Engine with no middlewares
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}

	// engine is the top group of gingle
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

// Default returns a instance of Engine with common middlewares
func Default() *Engine {
	engine := New()

	// Apply Logger and Recovery to engine
	engine.Use(Logger(), Recovery())

	return engine
}

// Run encapsulates ListenAndServe by using engine as handler
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP defines how gingle handle requests and responses
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	middlewares := make([]HandlerFunc, 0)
	for _, group := range engine.groups {
		// If belongs to some group, retrieve its middlewares
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	// Contruct the context
	// recording response writer, request and middlewares
	ctx := newContext(rw, req)
	ctx.middlewares = middlewares
	ctx.engine = engine

	// Parse the context
	// and handle in static or dynamic mode
	engine.router.handle(ctx)
}

// SetFuncMap sets render function
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob excutes template render
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.templates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
