package tcp_test

import (
	"testing"

	"github.com/rvflash/tcp"
)

func ExampleNew() {
	srv := tcp.New()
	srv.SYN(sleep)
	// now runs it!
}

func TestDefault(t *testing.T) {
	srv := tcp.Default()
	srv.ACK(stumble)
	// now runs it!
}
