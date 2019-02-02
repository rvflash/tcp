package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
)

// Conn is implemented by any Context.
type Conn interface {
	io.WriteCloser
	// Closed listening the cancellation context.
	Closed() <-chan struct{}
	// RawData returns the request's message.
	RawData() []byte
	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
}

// Context ...
type Context struct {
	cancel context.CancelFunc
	conn   net.Conn
	ctx    context.Context
	msg    []byte
	srv    *Server
}

func (c *Context) closing() {
	for _, f := range c.srv.out {
		f(c)
	}
	if err := c.Close(); err != nil {
		c.srv.errorf("context closing failed with %q", err)
	}
}

func (c *Context) incoming() {
	for _, f := range c.srv.in {
		f(c)
	}
}

func (c *Context) listening(d []byte) {
	for _, f := range c.srv.msg {
		f(c.copy(d))
	}
}

func (c *Context) copy(d []byte) *Context {
	var cc = *c
	cc.msg = make([]byte, len(d))
	copy(cc.msg, d)
	return &cc
}

func (c *Context) handle(ctx context.Context) {
	// Initiates the connection with a context by cancellation.
	c.ctx, c.cancel = context.WithCancel(ctx)
	// Launches any handler waiting for new connection.
	c.incoming()
	r := bufio.NewReader(c.conn)
	for {
		// For each new message
		d, err := r.ReadBytes('\n')
		if err != nil {
			// Closes the connection and the context.
			c.closing()
			return
		}
		c.listening(d)
	}
}

// Close implements the Conn interface.
func (c *Context) Close() error {
	if c.conn == nil {
		return nil
	}
	c.cancel()
	return c.conn.Close()
}

// Closed implements the Conn interface.
func (c *Context) Closed() <-chan struct{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Done()
}

// Data implements the Conn interface.
func (c *Context) RawData() []byte {
	return c.msg
}

// RemoteAddr implements the Conn interface.
func (c *Context) RemoteAddr() net.Addr {
	if c.conn == nil {
		return nil
	}
	return c.conn.RemoteAddr()
}

// String writes the given string into the connection.
func (c *Context) String(s string) {
	_, err := c.Write([]byte(s + "\n"))
	if err != nil {
		// todo panic
		c.srv.errorf("failed to write: %s", err)
	}
}

// Write implements the Conn interface.
func (c *Context) Write(d []byte) (int, error) {
	return c.conn.Write(d)
}
