module gingle-web

go 1.13

require (
	gingle v0.0.0
	gingle/middlewares v0.0.0
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
)

replace gingle => ./gingle

replace gingle/middlewares => ./gingle/middlewares

// 从 go 1.11 版本开始，引用相对路径的 package 需要使用上述方式
// 使用 go build main.go，不用 go run main.go
