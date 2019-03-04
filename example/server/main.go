package main

import (
	"log"

	"github.com/rvflash/tcp"
)

func main() {
	r := tcp.Default()
	r.ACK(func(c *tcp.Context) {
		// new message received
		body, err := c.ReadAll()
		if err != nil {
			c.Error(err)
		}
		log.Println(string(body))
		c.String("read")
	})
	r.SYN(func(c *tcp.Context) {
		c.String("hello")
	})
	r.FIN(func(c *tcp.Context) {
		log.Println("bye")
	})
	log.Fatal(r.Run(":9090"))
}
