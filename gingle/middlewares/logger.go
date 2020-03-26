package middlewares

import (
	"gingle"
	"log"
	"time"
)

func Logger() gingle.HandlerFunc {
	return func(ctx *gingle.Context) {
		t := time.Now()
		ctx.Next()
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
