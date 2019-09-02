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

	// Waiting for messages
	r := bufio.NewReader(c.rwc)
	for {
		select {
		case <-ctx.Done():
			// Connection closing, stops serving.
			c.bySegment(ctx, FIN, r)
			return
		default:
		}
		d, err := r.ReadBytes('\n')
		r := bytes.NewReader(d)
		if err != nil {
			// Unable to read on it: closing the connection.
			c.bySegment(ctx, FIN, r)
			return
		}
		// new message received
		c.bySegment(ctx, ACK, r)
	}
}
