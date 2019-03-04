package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		log.Fatalf("conn: %s", err)
	}
	var (
		s string
		r *bufio.Reader
	)
	for {
		// reads from stdin
		log.Print("> ")
		r = bufio.NewReader(os.Stdin)
		s, err = r.ReadString('\n')
		if err != nil {
			log.Fatalf("read stdin: %s", err)
		}
		// sends it
		_, err := fmt.Fprintf(conn, s)
		if err != nil {
			log.Fatalf("write conn: %s", err)
		}
		// listens for reply
		s, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatalf("read conn: %s", err)
		}
		log.Print("< " + s)
	}
}
