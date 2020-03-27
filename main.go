package main

import (
	"gingle"
	"gingle/middlewares"
	"net/http"
)

func main() {
	router := gingle.New()

	router.Static("/assets", "./static")

	router.GET("/testHTML", func(ctx *gingle.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Gingle!</h1>")
	})

	testString := router.Group("/testString")
	{
		testString.GET("/", func(ctx *gingle.Context) {
			ctx.String(http.StatusOK, "Mode = Static\nMessage = %s\nPattern = %s\nMethod = %s\n", ctx.Query("msg"), ctx.Pattern, ctx.Method)
		})

		testString.GET("/:msg", func(ctx *gingle.Context) {
			ctx.String(http.StatusOK, "Mode = Dynamic\nMessage = %s\nPattern = %s\nMethod = %s\n", ctx.Param("msg"), ctx.Pattern, ctx.Method)
		})
	}
	testString.Use(middlewares.Logger())

	testJSON := router.Group("/testJSON")
	{
		testJSON.POST("/", func(ctx *gingle.Context) {
			ctx.JSON(http.StatusOK, gingle.H{
				"username": ctx.PostForm("username"),
				"password": ctx.PostForm("password"),
			})
		})

		testJSON.POST("/*info", func(ctx *gingle.Context) {
			ctx.JSON(http.StatusOK, gingle.H{
				"info": ctx.Param("info"),
			})
		})
	}

	router.Run(":8080")
}
