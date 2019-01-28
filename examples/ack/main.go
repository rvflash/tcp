package main

import (
	"log"

	"github.com/rvflash/tcp"
)

func main() {
	r := tcp.Default()
	r.ACK(func(c tcp.Conn) {
		// new message received
		log.Printf("%s\n", c.RawData())
	})
	log.Fatal(r.Run(":9090"))
}
