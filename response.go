package tcp

import (
	"io"
	"net"
)

const noWritten = 1

// ResponseWriter interface is used by an TCP handler to construct the response.
type ResponseWriter interface {
	// Writer is the interface that wraps the basic Write method.
	io.Writer
	// WriteString allow to directly write string.
	WriteString(s string) (n int, err error)
	// Size returns the number of bytes already written into the response body.
	// -1: not already written
	Size() int
}

type responseWriter struct {
	conn io.Writer
	size int
}

func newResponseWriter(c net.Conn) *responseWriter {
	return &responseWriter{
		conn: c,
		size: noWritten,
	}
}

// Write implements the ResponseWriter interface.
func (w *responseWriter) Size() int {
	return w.size
}

// Write implements the ResponseWriter interface.
func (w *responseWriter) Write(p []byte) (n int, err error) {
	n, err = w.conn.Write(p)
	w.incr(n)
	return
}

// Write implements the ResponseWriter interface.
func (w *responseWriter) WriteString(s string) (n int, err error) {
	n, err = io.WriteString(w.conn, s)
	w.incr(n)
	return
}

func (w *responseWriter) incr(n int) {
	if n == noWritten {
		n = 0
	}
	w.size += n
}

func (w *responseWriter) rebase(conn io.Writer) {
	w.conn = conn
	w.size = noWritten
}
