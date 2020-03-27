package gingle

import (
	"net/http"
	"path"
)

// TODO: HandlerFunc 业务处理函数
type HandlerFunc func(*Context)

// TODO: RouterGroup 分组路由嵌套
type RouterGroup struct {
	prefix      string
	parent      *RouterGroup
	middlewares []HandlerFunc
	engine      *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	subgroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	// TODO: 无所谓先后的顺序
	engine.groups = append(engine.groups, subgroup)
	return subgroup
}

func (group *RouterGroup) register(method string, subpattern string, handler HandlerFunc) {
	// TODO: 组合继承与匿名成员
	pattern := group.prefix + subpattern
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.register("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.register("POST", pattern, handler)
}

func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.register("PUT", pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.register("DELETE", pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs)) // TODO: 关键转换

	return func(ctx *Context) {
		// TODO: 检查是否存在
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Fail(http.StatusNotFound, "404 Not Found")
			return
		}

		// TODO: 进行渲染返回
		fileServer.ServeHTTP(ctx.RespWriter, ctx.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	pattern := path.Join(relativePath, "/*filepath")
	group.GET(pattern, handler)
}
