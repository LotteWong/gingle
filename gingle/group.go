package gingle

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
