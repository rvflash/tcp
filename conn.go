package tcp

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
)

type conn struct {
	addr string
	rwc  net.Conn
	srv  *Server
}

func (c *conn) bySegment(ctx context.Context, segment string, body io.Reader) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	w := newWriter(c.rwc)
	req := c.newRequest(segment, body).WithContext(ctx)
	c.srv.ServeTCP(w, req)
}

func (c *conn) newRequest(segment string, body io.Reader) *Request {
	req := NewRequest(segment, body)
	req.RemoteAddr = c.addr
	return req
}

func (c *conn) serve(ctx context.Context) {
	// New connection
	c.bySegment(ctx, SYN, nil)
	// Connection closed
	defer c.bySegment(ctx, FIN, nil)
	// Waiting for messages
	r := bufio.NewReader(c.rwc)
	for {
		cb := make(chan []byte, 1)
		go func() {
			d, err := r.ReadBytes('\n')
			if err != nil {
				return
			}
			cb <- d
		}()
		select {
		case <-ctx.Done():
			return
		case b := <-cb:
			c.bySegment(ctx, ACK, bytes.NewReader(b))
		}
	}
}
