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
	ctx  context.Context
	srv  *Server
	rwc  net.Conn
}

func (c *conn) bySegment(segment string, body io.Reader) {
	w := newWriter(c.rwc)
	req := c.newRequest(segment, body)
	c.srv.ServeTCP(w, req)
}

func (c *conn) newRequest(segment string, body io.Reader) *Request {
	req := NewRequest(segment, body)
	req.RemoteAddr = c.addr
	return req.WithCancel(c.ctx)
}

func (c *conn) serve() {
	// deals with a new connection
	go c.bySegment(SYN, nil)
	// waiting for messages
	r := bufio.NewReader(c.rwc)
	for {
		d, err := r.ReadBytes('\n')
		r := bytes.NewReader(d)
		if err != nil {
			// unable to read on it: closing the connection.
			c.bySegment(FIN, r)
			return
		}
		// new message received
		go c.bySegment(ACK, r)
	}
}
