package main

import (
	"gingle"
	"net/http"
)

func main() {
	router := gingle.Default()

	router.Static("/assets", "./static")

	testString := router.Group("/testString")
	{
		testString.GET("/", func(ctx *gingle.Context) {
			ctx.String(http.StatusOK, "Mode = Static\nMessage = %s\nPattern = %s\nMethod = %s\n", ctx.Query("msg"), ctx.Pattern, ctx.Method)
		})

		testString.GET("/:msg", func(ctx *gingle.Context) {
			ctx.String(http.StatusOK, "Mode = Dynamic\nMessage = %s\nPattern = %s\nMethod = %s\n", ctx.Param("msg"), ctx.Pattern, ctx.Method)
		})
	}

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

	router.GET("/panic", func(ctx *gingle.Context) {
		names := []string{"bot", "exp"}
		ctx.String(http.StatusOK, "%s", names[2])
	})

	router.Run(":8080")
}
