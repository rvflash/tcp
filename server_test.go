package tcp_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/matryer/is"

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

func TestServer_Any(t *testing.T) {
	var (
		any = handleResp(tcp.ANY, false)
		dt  = []struct {
			req *tcp.Request
			out string
		}{
			{req: tcp.NewRequest(tcp.ANY, nil), out: any + handleResp(tcp.SYN, true)},
			{req: tcp.NewRequest(tcp.ACK, nil), out: any + handleResp(tcp.ACK, true)},
			{req: tcp.NewRequest(tcp.FIN, nil), out: any + handleResp(tcp.FIN, true)},
			{req: tcp.NewRequest(tcp.SYN, nil), out: any + handleResp(tcp.SYN, true)},
		}
		are = is.New(t)
	)
	// Listens all segments
	srv := tcp.New()
	for _, seg := range []string{tcp.ANY, tcp.ACK, tcp.FIN, tcp.SYN} {
		srv.Any(seg, handle(seg))
	}
	// Launches the test cases
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			rec := tcp.NewRecorder()
			srv.ServeTCP(rec, tt.req)
			are.Equal(rec.Body.String(), tt.out)
		})
	}
}

func handle(segment string) tcp.HandlerFunc {
	return func(c *tcp.Context) {
		c.String(handleResp(segment, c.Request.Segment == segment))
	}
}

func handleResp(segment string, expected bool) string {
	return fmt.Sprintf("%q segment: %t\n", segment, expected)
}
