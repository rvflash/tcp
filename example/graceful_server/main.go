package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rvflash/tcp"
)

func main() {
	bye := make(chan os.Signal, 1)
	signal.Notify(bye, os.Interrupt, syscall.SIGTERM)

	r := tcp.Default()
	r.ReadTimeout = 20 * time.Second
	r.ACK(func(c *tcp.Context) {
		// new message received
		body, err := c.ReadAll()
		if err != nil {
			c.Error(err)
			return
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

	go func() {
		err := r.Run(":9090")
		if err != nil {
			log.Printf("server: %q\n", err)
		}
	}()

	<-bye
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err := r.Shutdown(ctx)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}
