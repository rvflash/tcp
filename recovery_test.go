package tcp_test

import (
	"testing"

	"github.com/rvflash/tcp"
)

func TestRecovery(t *testing.T) {
	srv := tcp.New()
	srv.Use(tcp.Recovery())
	req := tcp.NewRequest(tcp.SYN, nil)
	w := tcp.NewRecorder()
	srv.ServeTCP(w, req)
}
