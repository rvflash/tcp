package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
)

// Conn is implemented by any client.
type Conn interface {
	io.WriteCloser
	// Closed listening the cancellation context.
	Closed() <-chan struct{}
	// RawData returns the latest message received.
	RawData() []byte
}

type client struct {
	cancel context.CancelFunc
	conn   net.Conn
	ctx    context.Context
	msg    []byte
	srv    *Server
}

func (c *client) closing() {
	for _, f := range c.srv.out {
		f(c)
	}
	if err := c.Close(); err != nil {
		c.srv.errorf("context closing failed with %q", err)
	}
}

func (c *client) incoming() {
	for _, f := range c.srv.in {
		f(c)
	}
}

func (c *client) listening(d []byte) {
	for _, f := range c.srv.msg {
		f(c.copy(d))
	}
}

func (c *client) copy(d []byte) *client {
	var cc = *c
	cc.msg = make([]byte, len(d))
	copy(cc.msg, d)
	return &cc
}

func (c *client) handle(ctx context.Context) {
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
func (c *client) Close() error {
	if c.conn == nil {
		return nil
	}
	c.cancel()
	return c.conn.Close()
}

// Closed implements the Conn interface.
func (c *client) Closed() <-chan struct{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Done()
}

// Data implements the Conn interface.
func (c *client) RawData() []byte {
	return c.msg
}

// Write implements the Conn interface.
func (c *client) Write(d []byte) (int, error) {
	return c.conn.Write(d)
}
