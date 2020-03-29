# gingle-web

A simple gin-like web framework implemented by Golang.

---

## Features

- [x] Serve Static Files
- [x] Handle Route Mappings <u>(static and dynamic)</u>
- [x] Middleware Plugin <u>(default and custom)</u>
- [X] Group Control
- [x] Context Support
- [x] Template Render

## Quick Start

### Initiation

```go
// Get an empty engine with no middlewares
router := gingle.New()

// Get an default engine with Logger() and Recovery() middlewares
router := gingle.Default
```

### Serve Static Files

```go
// Convert relative path to server root
router.Static("relativePath", "root")
```

### Handle Route Mappings

```go
router.GET("/pattern", func(ctx *gingle.Context) {
  // Static Route Mappings in GET method
})

router.POST("/pattern", func(ctx *gingle.Context) {
  // Static Route Mappings in POST method
})

router.PUT("/pattern", func(ctx *gingle.Context) {
  // Static Route Mappings in PUT method
})

router.DELETE("/pattern", func(ctx *gingle.Context) {
  // Static Route Mappings in DELETE method
})

router.GET("/:pattern", func(ctx *gingle.Context) {
  // Dynamic Route Mappings in /:pattern mode

  // Example:
  // Note - parameter `language` will match value `en` as follows
  // Route - /:language/content/
  // URL - http:localhost:8080/en/content/
})

router.GET("/*pattern", func(ctx *gingle.Context) {
  // Dynamic Route Mappings in /*pattern mode

  // Example:
  // Note - parameter `filepath` will match value `static/index.html` as follows
  // Route - /*filepath
  // URL - http:localhost:8080/static/index.html
})
```

### Middleware Plugin

```go
// `mymiddleware.go`: Define a custom middleware
func MyMiddleware() HandlerFunc {
  return func(ctx *Context) {
    // Do something before excuting handler

    ctx.Next()

    // Do something after excuting handler
  }
}

// `main.go`: Apply the middleware
router.Use(MyMiddleware())
```

### Group Control

```go
// Define a router group
api := router.Group("/api")

// Add routes to the router group
{
  api.POST("/signin", func(ctx *Context) {
    // pattern is /api/signin
  })
  api.POST("/signup", func(ctx *Context) {
    // pattern is /api/signup
  })
}

// Apply middlewares to the router group
api.Use(Logger(), Recovery())
```

## TODOs

- [ ] Helpers for PUT and DELETE
- [ ] Regular Expression Support
- [ ] Pratical Middlewares
- [ ] Unit and Benchmark Tests
