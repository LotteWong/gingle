module gingle-web

go 1.13

// Try to use go build instead of go run
require (
	gingle v0.0.0
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
)

// From go 1.11 version, import packages with relative path
// by using replace `ailas` => `path`
replace gingle => ./gingle
