package main

import (
	"gingle"
	"net/http"
)

func main() {
	router := gingle.New()

	router.GET("/testString", func(ctx *gingle.Context) {
		ctx.String(http.StatusOK, "Message = %s\nPattern = %s\nMethod = %s\n", ctx.Query("msg"), ctx.Pattern, ctx.Method)
	})

	router.GET("/testString/:msg", func(ctx *gingle.Context) {
		ctx.String(http.StatusOK, "Message = %s\nPattern = %s\nMethod = %s\n", ctx.Param("msg"), ctx.Pattern, ctx.Method)
	})

	router.POST("/testJSON", func(ctx *gingle.Context) {
		ctx.JSON(http.StatusOK, gingle.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	router.POST("/testJSON/*info", func(ctx *gingle.Context) {
		ctx.JSON(http.StatusOK, gingle.H{
			"info": ctx.Param("info"),
		})
	})

	router.GET("/testHTML", func(ctx *gingle.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Gingle!</h1>")
	})

	router.Run(":8080")
}
