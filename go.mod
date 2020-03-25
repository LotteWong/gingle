module gingle-web

go 1.13

require gingle v0.0.0

replace gingle => ./gingle

// 从 go 1.11 版本开始，引用相对路径的 package 需要使用上述方式
// 使用 go build main.go，不用 go run main.go
