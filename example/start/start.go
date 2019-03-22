package main

import (
	"log"

	"github.com/rvflash/tcp"
)

func main() {
	// creates a server with a logger and a recover on panic as middlewares.
	r := tcp.Default()
	r.ACK(func(c *tcp.Context) {
		// new message received
		// gets the request body
		buf, err := c.ReadAll()
		if err != nil {
			c.Error(err)
			return
		}
		// writes something as response
		c.String(string(buf))
	})
	err := r.Run(":9090") // listen and serve on 0.0.0.0:9090
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
}
