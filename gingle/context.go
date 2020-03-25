package gingle

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: 响应的结构体
type H map[string]interface{}

type Context struct {
	RespWriter http.ResponseWriter
	Req *http.Request
	
	Pattern string
	Method string
	
	StatusCode int
}

func newContext(rw http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		RespWriter: rw,
		Req: req,
		Pattern: req.URL.Path,
		Method: req.Method,
	}
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) SetStatus(code int) {
	ctx.StatusCode = code
	ctx.RespWriter.WriteHeader(ctx.StatusCode) // TODO: 请求行
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.RespWriter.Header().Set(key, value) // TODO: 请求头
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.RespWriter.Write([]byte(fmt.Sprintf(format, values...))) // TODO: 请求体
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.SetStatus(code)
	ctx.RespWriter.Write(data) // TODO: 请求体
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "application/json")
	bytes, _ := json.Marshal(obj)
	ctx.RespWriter.Write(bytes) // TODO: 请求体
}

func (ctx *Context) HTML(code int, html string) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "text/html")
	ctx.RespWriter.Write([]byte(html)) // TODO: 请求体
}
