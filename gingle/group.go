package gingle

import (
	"net/http"
	"path"
)

// HandlerFunc is business logic handler
type HandlerFunc func(*Context)

// RouterGroup is the group of routes
type RouterGroup struct {
	prefix      string        // group identity
	parent      *RouterGroup  // parent group
	middlewares []HandlerFunc // applicable middlewares
	engine      *Engine       // powerful router
}

// Group returns a instance of RouterGroup identified with prefix
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	subgroup := &RouterGroup{
		prefix: group.prefix + prefix, // Support nesting groups
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, subgroup)
	return subgroup
}

// register concatenates prefix and pattern then add it to router
func (group *RouterGroup) register(method string, subpattern string, handler HandlerFunc) {
	pattern := path.Join(group.prefix + subpattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET method encapsulation
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.register("GET", pattern, handler)
}

// POST method encapsulation
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.register("POST", pattern, handler)
}

// PUT method encapsulation
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.register("PUT", pattern, handler)
}

// DELETE method encapsulation
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.register("DELETE", pattern, handler)
}

// Use applies middlewares to group in order
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// createStaticHandler map client path to server path
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)             // client path
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs)) // server path

	return func(ctx *Context) {
		// Check whether the file exists
		file := ctx.Param("filepath") // custom parameter name
		if _, err := fs.Open(file); err != nil {
			ctx.Fail(http.StatusNotFound, "Status Not Found: "+ctx.Pattern)
			return
		}

		// Serve file content
		fileServer.ServeHTTP(ctx.RespWriter, ctx.Req)
	}
}

// Static registers a specific GET method to serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	pattern := path.Join(relativePath, "/*filepath") // custom parameter name
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	group.GET(pattern, handler)
}
