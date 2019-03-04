package tcp

import (
	"fmt"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				c.Error(&Error{
					msg:     "panic recovered",
					cause:   fmt.Errorf("%v", r),
					recover: true,
				})
				c.Abort()
			}
		}()
		// Processes the request
		c.Next()
	}
}
