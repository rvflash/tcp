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
	go c.bySegment(ctx, SYN, nil)
	// Waiting for messages
	r := bufio.NewReader(c.rwc)
	for {
		d, err := r.ReadBytes('\n')
		if err != nil {
			break
		}
		go c.bySegment(ctx, ACK, bytes.NewReader(d))
	}
	// Connection closed
	c.bySegment(ctx, FIN, nil)
}
