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

// ServeTCP implements the Handler interface.
func (c *conn) ServeTCP(w ResponseWriter, req *Request) {
	ctx := c.srv.get()
	ctx.writer.rebase(w)
	ctx.Request = req
	ctx.reset()
	c.handle(ctx)
	c.srv.put(ctx)
}

func (c *conn) bySegment(segment string, body io.Reader) {
	req := c.newRequest(segment, body)
	w := c.newResponseWriter()
	c.ServeTCP(w, req)
}

func (c *conn) handle(ctx *Context) {
	ctx.handlers = c.srv.computeHandlers(ctx.Request.Segment)
	if len(ctx.handlers) == 0 {
		return
	}
	ctx.Next()
}

func (c *conn) newResponseWriter() *responseWriter {
	return &responseWriter{
		ResponseWriter: c.rwc,
		size:           noWritten,
	}
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
