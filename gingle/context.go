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
	Req        *http.Request

	Pattern string
	Method  string
	Params  map[string]string

	StatusCode int

	middlewares []HandlerFunc
	index       int
}

func newContext(rw http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		RespWriter: rw,
		Req:        req,
		Pattern:    req.URL.Path,
		Method:     req.Method,
		index:      -1,
	}
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Param(key string) string {
	return ctx.Params[key]
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

func (ctx *Context) Fail(code int, err string) {
	ctx.index = len(ctx.middlewares) // TODO: 置为末位表示出错
	ctx.JSON(code, H{
		"msg": err,
	})
}

func (ctx *Context) Next() {
	// TODO: 正序处理
	ctx.index++

	for ; ctx.index < len(ctx.middlewares); ctx.index++ {
		ctx.middlewares[ctx.index](ctx)
	}
}
