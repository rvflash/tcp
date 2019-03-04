package tcp_test

import (
	"testing"

	"github.com/rvflash/tcp"
)

func TestRecovery(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("expected no panic")
		}
	}()
	srv := tcp.New()
	srv.Use(tcp.Recovery())
	srv.SYN(oops)

	w := tcp.NewRecorder()
	srv.ServeTCP(w, tcp.NewRequest(tcp.SYN, nil))
}

func TestRecovery2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	srv := tcp.New()
	srv.SYN(oops)

	w := tcp.NewRecorder()
	srv.ServeTCP(w, tcp.NewRequest(tcp.SYN, nil))
}

func oops(c *tcp.Context) {
	panic("oops, sorry!")
}
