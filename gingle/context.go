package gingle

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H helps convert obj to json
type H map[string]interface{}

// Context maintains information about each connect
type Context struct {
	// For managing responses and requests
	RespWriter http.ResponseWriter
	Req        *http.Request

	// For parsing the request
	Pattern string
	Method  string
	Params  map[string]string

	// For writing the response
	StatusCode int

	// For monitoring middlewares
	middlewares []HandlerFunc
	index       int

	// For rendering templates
	engine *Engine
}

// newContext returns a instance of Context
func newContext(rw http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		RespWriter: rw,
		Req:        req,
		Pattern:    req.URL.Path,
		Method:     req.Method,
		index:      -1,
	}
}

// Query is for GET Method to retrieve data
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key) // Get from request
}

// PostForm is for POST Method to retrieve data
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key) // Get from request
}

// Param is for /: or /* to retrieve data
func (ctx *Context) Param(key string) string {
	return ctx.Params[key] // Get from context
}

// SetStatus sets the response status code
func (ctx *Context) SetStatus(code int) {
	ctx.StatusCode = code
	ctx.RespWriter.WriteHeader(ctx.StatusCode)
}

// SetHeader sets the response header kvpair
func (ctx *Context) SetHeader(key string, value string) {
	ctx.RespWriter.Header().Set(key, value)
}

// String quickly contructs response of text/plain content type
func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.RespWriter.Write([]byte(fmt.Sprintf(format, values...)))
}

// Data quickly contructs response of octet-stream content type
func (ctx *Context) Data(code int, data []byte) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "application/octet-stream")
	ctx.RespWriter.Write(data)
}

// JSON quickly contructs response of application/json content type
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "application/json")
	bytes, _ := json.Marshal(obj) // Serialize the content
	ctx.RespWriter.Write(bytes)
}

// HTML quickly contructs response of text/html content type
func (ctx *Context) HTML(code int, filename string, data interface{}) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "text/html")
	ctx.engine.templates.ExecuteTemplate(ctx.RespWriter, filename, data)
}

// Fail quickly contructs response of errors
func (ctx *Context) Fail(code int, desc string) {
	// Once failed, stop excuting the rest middlewares
	ctx.index = len(ctx.middlewares)

	// Support custom status code and error message
	ctx.JSON(code, H{
		"msg": desc,
	})
}

// Next controls middleware excution order
func (ctx *Context) Next() {
	ctx.index++ // Mark current middleware

	for ctx.index < len(ctx.middlewares) {
		ctx.middlewares[ctx.index](ctx) // Call next middleware
		ctx.index++                     // To exit the loop
	}
}
