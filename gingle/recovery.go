package gingle

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Recovery simply recovers panics and prints errors
func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				// Inform server of error
				log.Printf("%s\n\n", trace(err))
				// Inform client of error
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		ctx.Next()
	}
}

// trace records the stack information of error
func trace(err interface{}) string {
	var pcs [32]uintptr
	// Skip `Callers`, `trace` and `defer func`
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(fmt.Sprintf("%s\n\nTraceback:", err)) // error description
	for _, pc := range pcs[:n] {
		function := runtime.FuncForPC(pc)
		file, line := function.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line)) // stack description
	}
	return str.String()
}
