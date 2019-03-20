package tcp

import (
	"io"
	"io/ioutil"
	"strings"
	"time"
)

// M is a shortcut for map[string]interface{}
type M map[string]interface{}

// Context allows us to pass variables between middleware and manage the flow.
type Context struct {
	// Request contains information about the TCP request.
	Request *Request
	// ResponseWriter writes the response on the connection.
	ResponseWriter
	// Keys is a key/value pair allows data sharing inside the context of each request.
	Shared M

	errs     Errors
	index    int
	handlers []HandlerFunc
	srv      *Server
	writer   responseWriter
}

const abortIndex = 63

// Abort prevents pending handlers from being called,
// but not interrupt the current handler.
func (c *Context) Abort() {
	c.index = abortIndex
}

// Canceled is a shortcut to listen the request's cancellation.
func (c *Context) Canceled() <-chan struct{} {
	if c.Request == nil {
		return nil
	}
	return c.Request.Canceled()
}

// Close immediately closes the connection.
// An error is returned when we fail to do it.
func (c *Context) Close() error {
	return c.writer.Close()
}

// Error reports a new error.
func (c *Context) Error(err error) {
	c.errs = append(c.errs, err)
}

// Err explains what failed during the request.
// The method name is inspired by the context package.
func (c *Context) Err() Errors {
	return c.errs
}

// Get retrieves the value associated to the given key inside the embed shared memory.
// If it's not exists, the request's context value is used as fail over.
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Shared[key]
	if exists || c.Request == nil {
		return
	}
	// fail over based on context values
	value = c.Request.Context().Value(key)
	exists = value != nil
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (value bool) {
	v, ok := c.Get(key)
	if ok {
		value, _ = v.(bool)
	}
	return
}

// GetDuration returns the value associated with the key as a time duration.
func (c *Context) GetDuration(key string) (value time.Duration) {
	v, ok := c.Get(key)
	if ok {
		value, _ = v.(time.Duration)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (value float64) {
	v, ok := c.Get(key)
	if ok {
		value, _ = v.(float64)
	}
	return
}

// GetInt returns the value associated with the key as a int.
func (c *Context) GetInt(key string) (value int) {
	v, ok := c.Get(key)
	if ok {
		value, _ = v.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as a int64.
func (c *Context) GetInt64(key string) (value int64) {
	v, ok := c.Get(key)
	if ok {
		value, _ = v.(int64)
	}
	return
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (value string) {
	v, ok := c.Get(key)
	if ok {
		value, _ = v.(string)
	}
	return
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

// ReadAll return the stream data.
func (c *Context) ReadAll() ([]byte, error) {
	if c.Request == nil {
		return nil, ErrRequest
	}
	if c.Request.Body == nil {
		return nil, io.EOF
	}
	return ioutil.ReadAll(c.Request.Body)
}

// String writes the given string on the current connection.
func (c *Context) String(s string) {
	const eom = "\n"
	if !strings.HasSuffix(s, eom) {
		// sends it now, ending the message.
		s += eom
	}
	_, err := c.writer.WriteString(s)
	if err != nil {
		c.Error(err)
	}
}

// Write implements the Conn interface.
func (c *Context) Write(d []byte) (int, error) {
	return c.writer.Write(d)
}

/*
todo
func (c *Context) isAborted() bool {
	return c.index >= abortIndex
}
*/

func (c *Context) reset() {
	c.ResponseWriter = &c.writer
	c.Shared = make(M)
	c.handlers = nil
	c.index = -1
	c.errs = nil
}
