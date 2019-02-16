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
			log.Println("err:", err)
		}
		log.Println("request:", body)
		log.Println("say hi!")
		c.String("hi!")
	})
	r.SYN(func(c *tcp.Context) {
		log.Println("say hello!")
		c.String("hello")
	})
	r.FIN(func(c *tcp.Context) {
		log.Println("remote addr")
		log.Println(c.Request.RemoteAddr)
	})
	log.Fatal(r.Run(":9090"))
}
