package main

import (
	"log"

	"github.com/rvflash/tcp"
)

func main() {
	r := tcp.Default()
	r.ACK(func(c *tcp.Context) {
		// new message received
		log.Printf("%s", c.RawData())
		c.String("hi!")
	})
	r.SYN(func(c *tcp.Context) {
		c.String("hello")
	})
	r.FIN(func(c *tcp.Context) {
		log.Print(c.RemoteAddr())
	})
	log.Fatal(r.Run(":9090"))
}
