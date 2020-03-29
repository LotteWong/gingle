package gingle

import (
	"net/http"
	"strings"
)

// router maintains information about trie tree and handlers
type router struct {
	roots    map[string]*node       // trie tree map
	handlers map[string]HandlerFunc // handlers map
}

// newRouter returns a instance of router
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// convertPatternToParts helps split pattern to parts
func convertPatternToParts(pattern string) []string {
	items := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range items {
		// It does not matter whether the pattern ends with `/` or not
		// For example, /cn/content is equal to /cn/content/
		// parts = ["cn", "part"]
		if item != "" {
			parts = append(parts, item)
			// Once item has symbol `*`, force to break the loop
			// For example, /*filename/filehash only matches /*filename
			// parts = ["*filename"]
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute insert nodes into trie tree and add hanlder to map
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	parts := convertPatternToParts(pattern)

	// Request method serve as trie tree root
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	// Trie tree insertion
	r.roots[method].insert(pattern, parts, 0)

	// Handlers mapping
	r.handlers[key] = handler
}

// getRoute search the node and params by pattern from trie tree
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := convertPatternToParts(pattern) // concrete pattern
	params := make(map[string]string)

	// First check whether request method is valid
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		foundParts := convertPatternToParts(n.pattern) // concrete or wild pattern
		for index, part := range foundParts {
			// If `:` mode, replace wild pattern to the exact concrete param
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}

			// If `*` mode, replace wild pattern to the joined concrete params
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// getRoutes search all nodes by traversal from trie tree
func (r *router) getRoutes(method string) []*node {
	// First check whether request method is valid
	root, ok := r.roots[method]
	if !ok {
		return nil
	}

	// Collects all nodes in level order
	nodes := make([]*node, 0)
	root.traverse(&nodes)
	return nodes
}

// handle is to map pattern to handler in static or dynamic mode
func (r *router) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Pattern)

	if n != nil { // Found it, start excution
		// Store wild params
		ctx.Params = params
		// Use concrete or wild pattern (node pattern in trie tree)
		// instead of concrete pattern (context pattern from request)
		key := ctx.Method + "-" + n.pattern
		// Excute handler last
		ctx.middlewares = append(ctx.middlewares, r.handlers[key])
	} else { // Not found, report error
		ctx.Fail(http.StatusNotFound, "Status Not Found: "+ctx.Pattern)
	}

	// Switch to context to control excution order
	ctx.Next()
}
