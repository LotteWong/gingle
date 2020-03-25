package main

import (
	"fmt"
	"gingle"
	"net/http"
)

func main() {
	router := gingle.New()

	router.GET("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "req.URL.Path = %q\n", req.URL.Path)
	})

	router.GET("/test", func(rw http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(rw, "Header[%q] = %q\n", k, v)
		}
	})

	router.Run(":8080")
}
