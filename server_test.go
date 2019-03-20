package tcp_test

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rvflash/tcp"
)

func ExampleNew() {
	// Runs a server without any middleware, just a handler named sleep,
	// waiting for new connection.
	srv := tcp.New()
	srv.SYN(sleep)
	// error is ignored for the demo.
	_ = srv.Run(":9009")
}

func ExampleDefault() {
	// Runs a server with the default middlewares: logger and recover.
	// The stumble handler waiting for new message.
	srv := tcp.Default()
	srv.ACK(stumble)
	// error is ignored for the demo.
	_ = srv.Run(":9999")
}

func TestNew(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	srv := tcp.New()
	srv.ACK(oops)
	srv.ServeTCP(tcp.NewRecorder(), tcp.NewRequest(tcp.ACK, nil))
}

func TestDefault(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("expected no panic")
		}
	}()
	srv := tcp.Default()
	srv.ACK(oops)
	srv.ServeTCP(tcp.NewRecorder(), tcp.NewRequest(tcp.ACK, nil))
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
		srv = tcp.New()
	)
	// Listens all segments
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

const (
	eol         = "\n"
	clientAddr  = ":9123"
	hiMsg       = "hi, there's someone?" + eol
	receivedMsg = "received: %d bytes" + eol
	welcomeMsg  = "welcome" + eol
)

func TestServer_Run(t *testing.T) {
	// Prepares the server
	are := is.New(t)
	srv := tcp.New()
	srv.ACK(acknowledge)
	srv.SYN(welcome)
	srv.FIN(bye)
	go func() {
		err := srv.Run(clientAddr)
		are.NoErr(err)
	}()
	time.Sleep(time.Millisecond * 50)

	// Initiates the client
	cli, err := net.Dial("tcp", clientAddr)
	are.NoErr(err)
	defer func() {
		are.NoErr(cli.Close())
	}()

	// Welcome ?
	out, err := readConn(cli, len(welcomeMsg))
	are.NoErr(err)
	are.Equal(string(out), welcomeMsg)

	// Says hi!
	are.NoErr(writeConn(cli, hiMsg))
	out, err = readConn(cli, len(receivedMsg))
	are.NoErr(err)
	are.Equal(string(out), fmt.Sprintf(receivedMsg, len(hiMsg)))
}

func readConn(c io.Reader, size int) (out []byte, err error) {
	out = make([]byte, size)
	_, err = c.Read(out)
	return
}

func writeConn(w io.Writer, data string) (err error) {
	_, err = w.Write([]byte(data))
	return
}

func acknowledge(c *tcp.Context) {
	b, err := c.ReadAll()
	if err != nil {
		return
	}
	c.String(fmt.Sprintf(receivedMsg, len(b)))
}

func bye(_ *tcp.Context) {
	// do nothing
}

func welcome(c *tcp.Context) {
	c.String(welcomeMsg)
}
