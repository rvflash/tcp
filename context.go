package tcp

import (
	"bufio"
	"net"
	"strings"
)

// Context ...
type Context struct {
	Request *Request
	ResponseWriter

	conn     net.Conn
	errs     Errors
	index    int
	handlers []HandlerFunc
	writer   *responseWriter
}

// Close implements the Conn interface.
func (c *Context) Close() error {
	if c.Request != nil {
		c.Request.Cancel()
	}
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Closed implements the Conn interface.
func (c *Context) Closed() <-chan struct{} {
	if c.Request == nil {
		return nil
	}
	return c.Request.Closed()
}

// Error reports a new error.
func (c *Context) Error(err error) {
	c.errs = append(c.errs, err)
}

// Err explains what failed during the request.
// The method name is inspired of the context package.
func (c *Context) Err() error {
	return c.errs
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

// String writes the given string into the connection.
func (c *Context) String(s string) {
	if !strings.HasSuffix(s, "\n") {
		// sends it now
		s += "\n"
	}
	_, err := c.ResponseWriter.WriteString(s)
	if err != nil {
		c.Error(err)
	}
}

// Write implements the Conn interface.
func (c *Context) Write(d []byte) (int, error) {
	return c.ResponseWriter.Write(d)
}

func (c *Context) catch() {
	// Launches any handler waiting for new connection.
	c.applyHandlers(SYN)
	r := bufio.NewReader(c.conn)
	for {
		x := c.copy()
		d, err := r.ReadBytes('\n')
		if err != nil {
			x.close()
			return
		}
		go x.handle(d)
	}
}

func (c *Context) close() {
	// launches any handler waiting for closed connection.
	c.applyHandlers(FIN)
	// tries to close the connection and the context
	if err := c.Close(); err != nil {
		c.Error(NewError("close", err))
	}
}

func (c *Context) copy() *Context {
	var cc = *c
	cc.handlers = nil
	return &cc
}

func (c *Context) handle(d []byte) {
	// use response as entry point
	// launches any handler waiting for new message.
	c.applyHandlers(ACK)
}

func (c *Context) reset() {
	c.ResponseWriter = c.writer
	c.handlers = nil
	c.index = -1
	c.errs = c.errs[0:0]
}
